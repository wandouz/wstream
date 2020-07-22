package stream

func (s *WindowedStream) Reduce(reduceFunc functions.WindowReduce) *DataStream {
	stream := s.toDataStream()
	operator := s.operator.(WindowOperator)
	operator.SetReduceFunc(reduceFunc)
	s.combine(stream)
	return stream
}
