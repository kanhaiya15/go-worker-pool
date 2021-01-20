package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/kanhaiya15/go-worker-pool/consumer"
	"github.com/kanhaiya15/go-worker-pool/producer"
)

var (
	workerPoolSize   = 5
	producerPoolSize = 2
)

func init() {
	inWorkerPoolSize, err := strconv.Atoi(os.Getenv("WORKER_POOL_SIZE"))
	if err != nil {
		fmt.Println("error on getting env workerPoolSize : ", err.Error())
	} else {
		workerPoolSize = inWorkerPoolSize
	}

	inProducerPoolSize, err := strconv.Atoi(os.Getenv("PRODUCER_POOL_SIZE"))
	if err != nil {
		fmt.Println("error on getting env inProducerPoolSize : ", err.Error())
	} else {
		producerPoolSize = inProducerPoolSize
	}
}

func main() {
	fmt.Println("workerPoolSize : ( ", workerPoolSize, " ) producerPoolSize : ( ", producerPoolSize, " )")
	consumer := consumer.Consumer{
		IntermediateChan: make(chan uint64, 1),
		JobsChan:         make(chan uint64, workerPoolSize),
	}

	// Set up cancellation context and waitgroup
	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// Start producers and Add [producerPoolSize] to WaitGroup
	wg.Add(producerPoolSize)
	for i := 0; i < producerPoolSize; i++ {
		producer := producer.Producer{CallbackFunc: consumer.CallbackFunc}
		go producer.Start(ctx, wg, i)
	}

	// Start consumer with cancellation context passed
	go consumer.Start(ctx)

	// Start workers and Add [workerPoolSize] to WaitGroup
	wg.Add(workerPoolSize)
	for i := 0; i < workerPoolSize; i++ {
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
