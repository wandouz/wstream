package operator

import (
	"bytes"
	"container/list"
	"encoding/gob"
	"math"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wandouz/wstream/functions"
	"github.com/wandouz/wstream/intfs"
	"github.com/wandouz/wstream/runtime/operator/windowing"
	"github.com/wandouz/wstream/runtime/operator/windowing/assigners"
	"github.com/wandouz/wstream/runtime/operator/windowing/evictors"
	"github.com/wandouz/wstream/runtime/operator/windowing/triggers"
	"github.com/wandouz/wstream/runtime/operator/windowing/windows"
	"github.com/wandouz/wstream/types"
	"github.com/wandouz/wstream/utils"
)

// EvictWindow is evictable window operator
// assigner is due to assign windows for each record
// trigger is due to judge timing to emit window records
// evictor is due to evict records before or after window records are emitted
type EvictWindow struct {
	assigner assigners.WindowAssinger
	trigger  triggers.Trigger
	evictor  evictors.Evictor

	applyFunc  functions.ApplyFunc
	reduceFunc functions.ReduceFunc
	out        Emitter

	windowsGroup sync.Map
	// windowsGroup map[windowing.WindowID]*windowing.WindowCollection

	watermarkTime   time.Time
	eventTimer      *EventTimerService
	processingTimer *ProcessingTimerService
	triggerContext  *WindowTriggerContext
	assignerContext *WindowAssignerContext
}

// NewEvictWindow return evictable window object
// evictor is necessary
func NewEvictWindow(assigner assigners.WindowAssinger, trigger triggers.Trigger, evictor evictors.Evictor) *EvictWindow {
	if assigner == nil {
		assigner = assigners.NewGlobalWindow()
	}
	if trigger == nil {
		trigger = assigner.GetDefaultTrigger()
	}
	if evictor == nil {
		panic("EvictWindow must has an evictor")
	}
	w := &EvictWindow{
		assigner:     assigner,
		trigger:      trigger,
		evictor:      evictor,
		windowsGroup: sync.Map{},

		applyFunc: &byPassApplyFunc{},
	}
	w.processingTimer = NewProcessingTimerService(w, time.Second)
	w.eventTimer = NewEventTimerService(w)
	// bind this window to triggerContext factory
	w.triggerContext = NewWindowTriggerContext(w.processingTimer, w.eventTimer)
	w.assignerContext = NewWindowAssignerContext(w.processingTimer)
	return w
}

// New is a factory method to new an EvictWindow operator object
func (w *EvictWindow) New() intfs.Operator {
	window := NewEvictWindow(w.assigner, w.trigger, w.evictor)
	window.SetApplyFunc(w.newApplyFunc())
	window.SetReduceFunc(w.newReduceFunc())
	return window
}

func (w *EvictWindow) newApplyFunc() (udf functions.ApplyFunc) {
	encodedBytes := encodeFunction(w.applyFunc)
	reader := bytes.NewReader(encodedBytes)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&udf)
	if err != nil {
		panic(err)
	}
	return
}

func (w *EvictWindow) newReduceFunc() (udf functions.ReduceFunc) {
	if w.reduceFunc == nil {
		return
	}
	encodedBytes := encodeFunction(w.reduceFunc)
	reader := bytes.NewReader(encodedBytes)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&udf)
	if err != nil {
		panic(err)
	}
	return
}

func (w *EvictWindow) SetApplyFunc(f functions.ApplyFunc) {
	w.applyFunc = f
}

func (w *EvictWindow) SetReduceFunc(f functions.ReduceFunc) {
	w.reduceFunc = f
}

func (w *EvictWindow) handleRecord(record types.Record, out Emitter) {
	assignedWindows := w.assigner.AssignWindows(record, w.assignerContext)

	for _, window := range assignedWindows {
		if w.isWindowLate(window) {
			// drop window if it is event time and late
			logrus.Warnf("drop late window %+v for record %+v", window, record)
			continue
		}
		k := utils.HashSlice(record.Key())
		wid := windowing.NewWindowID(k, window)
		var coll *windowing.WindowCollection
		if c, ok := w.windowsGroup.Load(wid); ok {
			coll = c.(*windowing.WindowCollection)
			coll.Append(record)
		} else {
			// evict window not reduce records in collections, so don't pass reduceFunc into it
			coll = windowing.NewWindowCollection(window, record.Time(), record.Key(), nil)
			coll.Append(record)
			w.windowsGroup.Store(wid, coll)
		}
		ctx := w.triggerContext.New(wid, coll.Len())
		signal := w.trigger.OnItem(record, record.Time(), window, ctx)
		if signal.IsFire() {
			w.emitWindow(coll, out)
		}
		if signal.IsPurge() {
			coll.Dispose()
		}
		w.registerCleanupTimer(wid, window)
	}
}

func (w *EvictWindow) emitWindow(records *windowing.WindowCollection, out Emitter) {
	// for TimeEvictor records without timestamp is invalid
	// so for safty size should count only records with timestamp
	w.evictor.EvictBefore(records, int64(records.Len()))

	windowEmitter := NewWindowEmitter(records.Time(), records.Keys(), out)
	iterator := records.Iterator()
	if w.reduceFunc != nil {
		var acc types.Record
		for iter := iterator; iter != nil; iter = iter.Next() {
			if acc == nil {
				acc = iter.Value.(types.Record)
			} else {
				acc = w.reduceFunc.Reduce(acc, iter.Value.(types.Record))
			}
		}
		newl := list.New()
		if acc != nil {
			newl.PushBack(acc)
		}
		// TODO: encapsulation this
		w.applyFunc.Apply(newl.Front(), windowEmitter)
	} else {
		w.applyFunc.Apply(iterator, windowEmitter)
	}

	w.evictor.EvictAfter(records, int64(records.Len()))
}

func (w *EvictWindow) registerCleanupTimer(wid windowing.WindowID, window windows.Window) {
	if window.MaxTimestamp().Equal(time.Unix(math.MaxInt64, 0)) {
		return
	}
	if w.assigner.IsEventTime() {
		w.eventTimer.RegisterEventTimer(wid, window.MaxTimestamp())
	} else {
		w.processingTimer.RegisterProcessingTimer(wid, window.MaxTimestamp())
	}
}

func (w *EvictWindow) isWindowLate(window windows.Window) bool {
	return w.assigner.IsEventTime() && window.MaxTimestamp().Before(w.eventTimer.CurrentWatermarkTime())
}

// EvictWindow operator don't emit watermark from upstream operator
// and will emit new watermark when emit window
func (w *EvictWindow) handleWatermark(wm *types.Watermark, out Emitter) {
	w.eventTimer.Drive(wm.Time())
}

// onProcessingTime is callback for processing timer service
func (w *EvictWindow) onProcessingTime(wid windowing.WindowID, t time.Time) {
	c, ok := w.windowsGroup.Load(wid)
	if !ok {
		return
	}
	coll := c.(*windowing.WindowCollection)
	signal := w.trigger.OnProcessingTime(t, wid.Window())
	if signal.IsFire() {
		w.emitWindow(coll, w.out)
	}
	if signal.IsPurge() {
		coll.Dispose()
	}
	if !w.assigner.IsEventTime() && w.isCleanupTime(wid.Window(), t) {
		coll.Dispose()
		w.windowsGroup.Delete(wid)
	}
}

// onEventTIme is callback for event timer service
func (w *EvictWindow) onEventTime(wid windowing.WindowID, t time.Time) {
	c, ok := w.windowsGroup.Load(wid)
	if !ok {
		return
	}
	coll := c.(*windowing.WindowCollection)

	w.likelyEmitWatermark()

	signal := w.trigger.OnEventTime(t, wid.Window())
	if signal.IsFire() {
		w.emitWindow(coll, w.out)
	}
	if signal.IsPurge() {
		coll.Dispose()
	}
	if w.assigner.IsEventTime() && w.isCleanupTime(wid.Window(), t) {
		// clean window
		coll.Dispose()
		w.windowsGroup.Delete(wid)
	}
}

func (w *EvictWindow) isCleanupTime(window windows.Window, t time.Time) bool {
	return t.Equal(window.MaxTimestamp())
}

// check if should emit new watermark
func (w *EvictWindow) likelyEmitWatermark() {
	eventTime := w.eventTimer.CurrentWatermarkTime()
	if eventTime.After(w.watermarkTime) {
		w.out.Emit(types.NewWatermark(eventTime))
		w.watermarkTime = eventTime
	}
}

// Run this operator
func (w *EvictWindow) Run(in Receiver, out Emitter) {
	w.out = out
	w.processingTimer.Start()
	consume(in, out, w)
}
