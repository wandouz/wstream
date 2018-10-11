package assigners

import (
	"time"

	"github.com/wandouz/wstream/runtime/operator/windowing/triggers"
	"github.com/wandouz/wstream/runtime/operator/windowing/windows"
	"github.com/wandouz/wstream/types"
)

// SlidingProcessingTimeWindow assigner
// offset represent to timezone offset duration
type SlidingProcessingTimeWindow struct {
	period int64
	every  int64
	offset int64
}

func NewSlidingProcessingTimeWindow(period, every, offset int64) *SlidingProcessingTimeWindow {
	if offset < 0 || period <= 0 {
		panic("SlidingProcessingTimeWindow params must satisfy period > 0")
	}
	return &SlidingProcessingTimeWindow{
		period: period,
		every:  every,
		offset: offset,
	}
}

func (a *SlidingProcessingTimeWindow) AssignWindows(item types.Item, ctx AssignerContext) []windows.Window {
	var ret []windows.Window
	ts := ctx.GetCurrentProcessingTime().Unix()
	lastStart := GetWindowStartWithOffset(ts, a.offset, a.every)
	for start := lastStart; start > ts-a.period; start -= a.every {
		ret = append(ret, windows.NewTimeWindow(time.Unix(start, 0), time.Unix(start+a.period, 0)))
	}
	return ret
}

func (a *SlidingProcessingTimeWindow) GetDefaultTrigger() triggers.Trigger {
	return triggers.NewProcessingTimeTrigger()
}

func (a *SlidingProcessingTimeWindow) IsEventTime() bool {
	return false
}
