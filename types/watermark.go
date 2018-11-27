package types

//go:generate msgp -o codec_watermark_item.go

import (
	"time"
)

// Watermark is kind of item
type Watermark struct {
	T time.Time
}

func NewWatermark(t time.Time) *Watermark {
	return &Watermark{T: t}
}

func (wm *Watermark) Type() ItemType {
	return TypeWatermark
}

func (wm *Watermark) Clone() Item {
	return NewWatermark(wm.T)
}

func (wm *Watermark) Time() time.Time {
	return wm.T
}

func (wm *Watermark) SetTime(t time.Time) {
	wm.T = t
}

func (wm *Watermark) After(x *Watermark) bool {
	return wm.Time().After(x.Time())
}

func (wm *Watermark) AsRow() (Row, error) {
	encodedBytes, err := wm.MarshalMsg(nil)
	if err != nil {
		return Row{}, err
	}
	return Row{
		itemType: TypeWatermark,
		item:     encodedBytes,
	}, nil
}
