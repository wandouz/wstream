package stream

import (
	"github.com/wandouz/wstream/functions"
	"github.com/wandouz/wstream/runtime/operator"
)

func (s *KeyedStream) Reduce(reduceFunc functions.ReduceFunc) *DataStream {
	name := "reduce"
	graph := s.graph
	newStream := s.ToDataStream(name, nil)
	graph.AddStreamEdge(s, newStream)

	newStream.operator = operator.GenReduce(reduceFunc)
	return newStream
}
