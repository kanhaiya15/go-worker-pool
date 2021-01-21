package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kanhaiya15/go-worker-pool/cfg"
	"github.com/kanhaiya15/go-worker-pool/consumer"
	"github.com/kanhaiya15/go-worker-pool/producer"
)

func init() {
	cfg.Setup()
}

func main() {
	fmt.Println("ConsumerPoolSize : ( ", cfg.ConsumerPoolSize, " ) producerPoolSize : ( ", cfg.ProducerPoolSize, " )")
	consumer := consumer.Consumer{
		IntermediateChan: make(chan uint64, 1),
		JobsChan:         make(chan uint64, cfg.ConsumerPoolSize),
	}

	// Set up cancellation context and waitgroup
	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// Start producers and Add [producerPoolSize] to WaitGroup
	wg.Add(cfg.ProducerPoolSize)
	for i := 0; i < cfg.ProducerPoolSize; i++ {
		producer := producer.Producer{CallbackFunc: consumer.CallbackFunc}
		go producer.Start(ctx, wg, i)
	}

	// Start consumer with cancellation context passed
	go consumer.Start(ctx)

	// Start workers and Add [ConsumerPoolSize] to WaitGroup
	wg.Add(cfg.ConsumerPoolSize)
	for i := 0; i < cfg.ConsumerPoolSize; i++ {
		go consumer.WorkerFunc(wg, i)
	}

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan // Blocks here until interrupted

	// Handle shutdown
	fmt.Println("\n*********************************\nShutdown signal received\n*********************************")
	cancelFunc() // Signal cancellation to context.Context
	wg.Wait()    // Block here until are workers are done

	fmt.Println("All workers done, shutting down!")
}
