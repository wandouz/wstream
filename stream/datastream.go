package stream

import (
	"github.com/wandouz/wstream/intfs"
)

type DataStream struct {
	parallel int
	operator intfs.Operator

	// flow reference
	streamNode *StreamNode
	flow       *Flow
}

/*
DataStream API
*/

func NewDataStream(flow *Flow, parallel int) *DataStream {
	return &DataStream{
		parallel: parallel,
		flow:     flow,
	}
}

func (s *DataStream) Operator() intfs.Operator {
	return s.operator.New()
}

func (s *DataStream) Parallelism() int {
	return s.parallel
}

func (s *DataStream) clone() *DataStream {
	return &DataStream{
		flow:     s.flow,
		parallel: s.parallel,
	}
}

func (s *DataStream) SetPartition(parallel int) *DataStream {
	s.parallel = parallel
	return s
}

func (s *DataStream) SetStreamNode(node *StreamNode) {
	s.streamNode = node
}

func (s *DataStream) GetStreamNode() (node *StreamNode) {
	return s.streamNode
}

func (s *DataStream) toKeyedStream(keys ...interface{}) *KeyedStream {
	return NewKeyedStream(s.flow, s.parallel, keys)
}

func (s *DataStream) connect(stream Stream) {
	s.flow.AddStreamEdge(s, stream)
}
