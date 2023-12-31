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

func (eventBus *EventBus) AddTopic(topic string) {
	eventBus.Lock()
	eventBus.Unlock()

	eventBus.subscibers[topic] = make(subscribersChannels, 0)
}

func (eb *EventBus) Subsribe(topic string) eventChannel {
	eb.Lock()
	defer eb.Unlock()
	ec := make(eventChannel)

	if subscribers, found := eb.subscibers[topic]; found {
		eb.subscibers[topic] = append(subscribers, ec)
	} else {
		eb.subscibers[topic] = append(subscribersChannels{}, ec)
	}

	return ec
}
