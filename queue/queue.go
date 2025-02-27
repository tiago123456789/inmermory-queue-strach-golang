package queue

import (
	"sync"
	"time"
)

type QueueHandler func(message interface{}) error

type Consumer struct {
	Locked         chan bool
	MessageChannel chan interface{}
	Handler        QueueHandler
}

type IQueue interface {
	Publish(message interface{})
	AddHandler(callback QueueHandler)
	Start() *sync.WaitGroup
}

type Queue struct {
	mu         sync.Mutex
	data       []interface{}
	consumers  []Consumer
	dataFailed chan interface{}
}

func New() *Queue {
	return &Queue{
		data:       []interface{}{},
		dataFailed: make(chan interface{}),
	}
}

func (q *Queue) getMessage() interface{} {
	message := q.data[0]
	q.data = q.data[1:]
	return message
}

func (q *Queue) requeue(message interface{}) {
	q.data = append(q.data, message)
}

func (q *Queue) notifyConsumer(message interface{}) bool {
	hasConsumedMessage := false

	if len(q.consumers) == 0 {
		return hasConsumedMessage
	}

	for _, consumer := range q.consumers {
		if len(consumer.Locked) == cap(consumer.Locked) {
			continue
		} else {
			consumer.MessageChannel <- message
			return true
		}
	}

	return hasConsumedMessage
}

func (q *Queue) Publish(message interface{}) {
	hasConsumedMessage := q.notifyConsumer(message)
	if hasConsumedMessage == false {
		q.data = append(q.data, message)
	}
}

func (q *Queue) AddHandler(callback QueueHandler) {
	q.consumers = append(q.consumers, Consumer{
		Locked:         make(chan bool, 1),
		MessageChannel: make(chan interface{}, 1),
		Handler:        callback,
	})
}

func (q *Queue) Start() *sync.WaitGroup {
	var wg sync.WaitGroup

	for _, consumer := range q.consumers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case message := <-consumer.MessageChannel:
					consumer.Locked <- true
					err := consumer.Handler(message)
					if err != nil {
						q.dataFailed <- message
						time.Sleep(time.Second * 1)
					}
					<-consumer.Locked
				}
			}
		}()
	}

	go func() {
		for {
			select {
			case message := <-q.dataFailed:
				q.mu.Lock()
				q.requeue(message)
				q.mu.Unlock()
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				q.mu.Lock()
				if len(q.data) > 0 {
					message := q.getMessage()
					hasConsumed := q.notifyConsumer(message)
					if hasConsumed == false {
						q.requeue(message)
					}
				}

				q.mu.Unlock()
			}
		}
	}()

	return &wg
}
