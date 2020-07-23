package stream

import (
	"sync"

	"github.com/zhnpeng/wstream/runtime/execution"
	"github.com/zhnpeng/wstream/utils/graph"
)

/*
Flow is a DAG Graph organized with streams
*/
type Flow struct {
	Name        string
	Transformed bool
	Vertices    map[int]*StreamNode
	Graph       *graph.Mutable
}

func NewFlow(name string) *Flow {
	return &Flow{
		Name:     name,
		Vertices: make(map[int]*StreamNode),
		Graph:    graph.New(0),
	}
}

// StreamNode assign unique for each stream
// and combine stream with execution task
type StreamNode struct {
	ID     int
	Stream Stream
	task   *execution.Task
}

func newStreamNode(id int, s Stream) *StreamNode {
	return &StreamNode{
		ID:     id,
		Stream: s,
	}
}

func (f *Flow) GetStream(id int) (stm Stream) {
	if StreamNode, ok := f.Vertices[id]; ok {
		stm = StreamNode.Stream
	}
	return
}

func (f *Flow) GetTask(id int) (task *execution.Task) {
	if t, ok := f.Vertices[id]; ok {
		task = t.task
	}
	return
}

func (f *Flow) Len() int {
	return len(f.Vertices)
}

// AddStream add a stream vertex to Graph
func (f *Flow) AddStream(stm Stream) {
	id := f.Graph.AddVertex()
	StreamNode := newStreamNode(id, stm)
	stm.SetStreamNode(StreamNode)
	f.Vertices[id] = StreamNode
}

// AddStreamEdge add directed edge between two stream
func (f *Flow) AddStreamEdge(from, to Stream) error {
	if !f.existsStream(from) {
		f.AddStream(from)
	}
	if !f.existsStream(to) {
		f.AddStream(to)
	}
	fromID := from.GetStreamNode().ID
	toID := to.GetStreamNode().ID
	return f.Graph.AddEdge(fromID, toID)
}

// CombineStream join right strem to the left
// right stream will not add to Graph
// any StreamNode connect to right stream will collect to the left
func (f *Flow) CombineStream(left, right Stream) {
	if !f.existsStream(left) {
		f.AddStream(left)
	}
	right.SetStreamNode(left.GetStreamNode())
}

func (f *Flow) existsStream(stm Stream) bool {
	StreamNode := stm.GetStreamNode()
	if StreamNode == nil {
		return false
	}
	if _, ok := f.Vertices[StreamNode.ID]; ok {
		return true
	}
	return false
}

func (f *Flow) GetStreamNode(id int) (StreamNode *StreamNode) {
	return f.Vertices[id]
}

func (f *Flow) BFSBoth(v int, do func(v, w int, c int64)) {
	graph.BFSAll(f.Graph, v, do)
}

// Run in local mode
func (f *Flow) Run() {
	f.LocalTransform()
	var wg sync.WaitGroup
	start := f.GetStreamNode(0)
	wg.Add(1)
	go func() {
		defer wg.Done()
		start.task.Run()
	}()
	graph.BFSBoth(f.Graph, 0, func(v, w int, c int64) {
		task := f.GetStreamNode(w).task
		wg.Add(1)
		go func() {
			defer wg.Done()
			task.Run()
		}()
	})
	wg.Wait()
}

// RunDistrib run in distribution mode
// TODO: implement this
func (f *Flow) RunDistrib() {

}
