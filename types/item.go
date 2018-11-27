package types

import "time"

type Item interface {
	Type() ItemType
	Clone() Item
	AsRow() (Row, error)
	Time() time.Time
	SetTime(t time.Time)
}
