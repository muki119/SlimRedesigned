package Stream

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
	streams "v1/Services/Streams"

	"golang.org/x/sync/semaphore"

	"github.com/redis/go-redis/v9"
)

type EventHandler func(context.Context, any) error

// StreamsEventBus Heavily Opinionated Redis Stream Manager , meant to act as a layer of abstraction from the redis stream
// hopefully allowing for easier management of the stream
type StreamsEventBus struct {
	Connection    *redis.Client
	ConsumerGroup string
	ConsumerName  string
	streamTable   map[string]EventHandler
	Services      *streams.StreamServices
	Listening     atomic.Bool

	Timeout   time.Duration
	waitGroup *sync.WaitGroup // wait group for Listen process - This is used for graceful closure.
	ctx       context.Context

	maxConcurrentSem *semaphore.Weighted // The max goroutines that can be spawned for processing tasks.
}
type EventBusConfig struct {
	ConsumerName       string
	ConsumerGroup      string
	Services           *streams.StreamServices
	Connection         *redis.Client
	timeout            time.Duration
	maxConcurrentProcs int64
}

func (config *EventBusConfig) CreateNewStreamsEventBus() *StreamsEventBus {
	newStreamEventBus := NewStreamsEventBus(
		config.ConsumerName,
		config.ConsumerGroup,
		config.Services,
		config.Connection,
		config.timeout,
		config.maxConcurrentProcs,
	)
	return newStreamEventBus
}

// NewStreamsEventBus The consumer group will be the same for all streams.
func NewStreamsEventBus(consumerName string, consumerGroup string, services *streams.StreamServices, conn *redis.Client, timeout time.Duration, maxConcurrentProcs int64) *StreamsEventBus {
	newStreamEventBus := &StreamsEventBus{
		Connection:       conn,
		ConsumerGroup:    consumerGroup,
		ConsumerName:     consumerName,
		streamTable:      make(map[string]EventHandler),
		Services:         services,
		Timeout:          timeout,
		waitGroup:        &sync.WaitGroup{},
		ctx:              context.Background(),
		maxConcurrentSem: semaphore.NewWeighted(maxConcurrentProcs),
	}
	newStreamEventBus.Listening.Store(false)
	return newStreamEventBus

}

// Initialize Joins all streams and consumer groups , and removes some pending messages in all streams associated
func (eventBus *StreamsEventBus) Initialize() error {
	for stream := range eventBus.streamTable { // create all the groups for all the streams
		_, err := eventBus.Connection.XGroupCreateMkStream(eventBus.ctx, stream, eventBus.ConsumerGroup, "$").Result()
		if err != nil && err.Error() != "BUSYGROUP" {
			return err
		}
	}
	err := eventBus.processPendingMessages()
	if err != nil {
		slog.Error("Error while processing pending events", err)
		return err
	}
	return nil
}

// StreamHandler Registers a handlerFunc for a stream s
//
//	All handler functions should be declared before init and listening
func (eventBus *StreamsEventBus) StreamHandler(stream string, handlerFunc EventHandler) {
	eventBus.streamTable[stream] = handlerFunc
}

func (eventBus *StreamsEventBus) Close() error {
	eventBus.Listening.Store(false)
	eventBus.waitGroup.Wait()
	err := eventBus.Connection.Close()
	if err != nil {
		return err
	}
	return nil
}
