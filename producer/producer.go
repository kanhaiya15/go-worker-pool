package producer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kanhaiya15/go-worker-pool/producer-consumer/utils"
)

// Producer simulates an external library that invokes the
// registered callback when it has new data for us once per 100ms.
type Producer struct {
	CallbackFunc func(event uint64)
}

// Start invokes the registered callback when it has new data for us once per 100ms.
func (p Producer) Start(ctx context.Context, wg *sync.WaitGroup, index int) {
	defer fmt.Println("Producer closing done!! ", index)
	defer wg.Done()
	nextUint64 := utils.NextUint64()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Producer received cancellation signal, closing Producer ", index)
			return
		default:
			i := nextUint64()
			fmt.Printf("Producer %d produced %d\n", index, i)
			p.CallbackFunc(i)
			time.Sleep(time.Millisecond * 100)
		}
	}
}
