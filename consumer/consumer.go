package consumer

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Consumer below here!
type Consumer struct {
	IntermediateChan chan uint64
	JobsChan         chan uint64
}

// CallbackFunc is invoked each time the external lib passes an event to us.
func (c Consumer) CallbackFunc(event uint64) {
	c.IntermediateChan <- event
}

// WorkerFunc starts a single worker function that will range on the JobsChan until that channel closes.
func (c Consumer) WorkerFunc(wg *sync.WaitGroup, index int) {
	defer wg.Done()

	fmt.Printf("Worker %d starting\n", index)
	for job := range c.JobsChan {
		// simulate work  taking between 1-3 seconds
		fmt.Printf("Worker %d started job %d\n", index, job)
		time.Sleep(time.Millisecond * time.Duration(1000+rand.Intn(2000)))
		fmt.Printf("Worker %d finished processing job %d\n", index, job)
	}
	fmt.Printf("Worker %d interrupted\n", index)
}

// Start acts as the proxy between the IntermediateChan and JobsChan, with a select to support graceful shutdown.
func (c Consumer) Start(ctx context.Context) {
	for {
		select {
		case job := <-c.IntermediateChan:
			c.JobsChan <- job
		case <-ctx.Done():
			fmt.Println("Consumer received cancellation signal, closing JobsChan!")
			close(c.JobsChan)
			fmt.Println("Consumer closed JobsChan")
			return
		}
	}
}
