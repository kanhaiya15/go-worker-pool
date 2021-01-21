package cfg

import (
	"fmt"
	"os"
	"strconv"
)

// Default
var (
	ConsumerPoolSize = 2
	ProducerPoolSize = 2
)

// Setup cfg
func Setup() {
	consumerPoolSize, err := strconv.Atoi(os.Getenv("CONSUMER_POOL_SIZE"))
	if err != nil {
		fmt.Println("error in getting ConsumerPoolSize from env, setting up default : ", ConsumerPoolSize)
	} else {
		ConsumerPoolSize = consumerPoolSize
	}

	producerPoolSize, err := strconv.Atoi(os.Getenv("PRODUCER_POOL_SIZE"))
	if err != nil {
		fmt.Println("error in getting ProducerPoolSize from env, setting up default : ", ProducerPoolSize)
	} else {
		ProducerPoolSize = producerPoolSize
	}
}
