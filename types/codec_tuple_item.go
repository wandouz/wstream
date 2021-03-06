package types

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *TupleRecord) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "T":
			z.T, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "K":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.K) >= int(zb0002) {
				z.K = (z.K)[:zb0002]
			} else {
				z.K = make([]interface{}, zb0002)
			}
			for za0001 := range z.K {
				z.K[za0001], err = dc.ReadIntf()
				if err != nil {
					return
				}
			}
		case "V":
			var zb0003 uint32
			zb0003, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.V) >= int(zb0003) {
				z.V = (z.V)[:zb0003]
			} else {
				z.V = make([]interface{}, zb0003)
			}
			for za0002 := range z.V {
				z.V[za0002], err = dc.ReadIntf()
				if err != nil {
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *TupleRecord) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "T"
	err = en.Append(0x83, 0xa1, 0x54)
	if err != nil {
		return
	}
	err = en.WriteTime(z.T)
	if err != nil {
		return
	}
	// write "K"
	err = en.Append(0xa1, 0x4b)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.K)))
	if err != nil {
		return
	}
	for za0001 := range z.K {
		err = en.WriteIntf(z.K[za0001])
		if err != nil {
			return
		}
	}
	// write "V"
	err = en.Append(0xa1, 0x56)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.V)))
	if err != nil {
		return
	}
	for za0002 := range z.V {
		err = en.WriteIntf(z.V[za0002])
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *TupleRecord) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "T"
	o = append(o, 0x83, 0xa1, 0x54)
	o = msgp.AppendTime(o, z.T)
	// string "K"
	o = append(o, 0xa1, 0x4b)
	o = msgp.AppendArrayHeader(o, uint32(len(z.K)))
	for za0001 := range z.K {
		o, err = msgp.AppendIntf(o, z.K[za0001])
		if err != nil {
			return
		}
	}
	// string "V"
	o = append(o, 0xa1, 0x56)
	o = msgp.AppendArrayHeader(o, uint32(len(z.V)))
	for za0002 := range z.V {
		o, err = msgp.AppendIntf(o, z.V[za0002])
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TupleRecord) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "T":
			z.T, bts, err = msgp.ReadTimeBytes(bts)
			if err != nil {
				return
			}
		case "K":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.K) >= int(zb0002) {
				z.K = (z.K)[:zb0002]
			} else {
				z.K = make([]interface{}, zb0002)
			}
			for za0001 := range z.K {
				z.K[za0001], bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					return
				}
			}
		case "V":
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.V) >= int(zb0003) {
				z.V = (z.V)[:zb0003]
			} else {
				z.V = make([]interface{}, zb0003)
			}
			for za0002 := range z.V {
				z.V[za0002], bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *TupleRecord) Msgsize() (s int) {
	s = 1 + 2 + msgp.TimeSize + 2 + msgp.ArrayHeaderSize
	for za0001 := range z.K {
		s += msgp.GuessSize(z.K[za0001])
	}
	s += 2 + msgp.ArrayHeaderSize
	for za0002 := range z.V {
		s += msgp.GuessSize(z.V[za0002])
	}
	return
}
