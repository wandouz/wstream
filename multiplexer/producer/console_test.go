package producer

import (
	"context"
	"testing"
	"time"

	"gopkg.in/tomb.v1"

	"github.com/zhnpeng/wstream/types"

	"github.com/zhnpeng/wstream/multiplexer"
)

func TestConsole_Produce(t *testing.T) {
	p := &Console{
		BasicProducer: &BasicProducer{
			Tomb:     &tomb.Tomb{},
			messages: make(chan multiplexer.Message, 10),
		},
		Format: "%v\n",
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	go p.Produce(ctx)

	for i := 0; i < 10; i++ {
		p.Write(multiplexer.Message{
			Data: types.NewRawMapRecord(map[string]interface{}{
				"value": i,
			}),
		})
	}

	time.Sleep(1 * time.Second)
	cancelFunc()
	err := p.Wait()
	if err != nil {
		t.Error(err)
	}
}
