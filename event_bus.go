package eventbus

import "sync"

type Event struct {
	Data  interface{}
	Topic string
}

type eventChannel chan Event

type subscribersChannels []eventChannel

type EventBus struct {
	subscibers map[string]subscribersChannels
	sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscibers: make(map[string]subscribersChannels),
	}
}

func (eventBus *EventBus) Publish(topic string, data interface{}) {
	event := Event{
		Data:  data,
		Topic: topic,
	}

	eventBus.RLock()
	channels := eventBus.subscibers[topic]
	eventBus.RUnlock()

	var wg sync.WaitGroup
	wg.Add(len(channels))

	go func(channels subscribersChannels, event Event) {
		for _, channel := range channels {
			channel <- event
			wg.Done()
		}
	}(channels, event)

	wg.Wait()
}

func (eventBus *EventBus) PublishAsync(topic string, data interface{}) {
	event := Event{
		Data:  data,
		Topic: topic,
	}

	eventBus.RLock()
	channels := eventBus.subscibers[topic]
	eventBus.RUnlock()

	go func(channels subscribersChannels, event Event) {
		for _, channel := range channels {
			channel <- event
		}
	}(channels, event)
}
