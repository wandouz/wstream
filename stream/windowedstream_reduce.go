package stream

import (
	"github.com/wandouz/wstream/functions"
)

func (s *WindowedStream) Reduce(reduceFunc functions.Reduce) *DataStream {
	stream := s.toDataStream()
	operator := s.operator.(WindowOperator)
	operator.SetReduceFunc(reduceFunc)
	s.leftMerge(stream)
	return stream
}
