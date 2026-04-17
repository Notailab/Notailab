package dao

import "sync"

type FileChangeEvent struct {
	Action    string `json:"action"`
	ProjectID uint   `json:"project_id"`
	FileID    uint   `json:"file_id,omitempty"`
}

type fileEventBroker struct {
	mu          sync.RWMutex
	subscribers map[uint]map[chan FileChangeEvent]struct{}
}

var globalFileEventBroker = newFileEventBroker()

func newFileEventBroker() *fileEventBroker {
	return &fileEventBroker{
		subscribers: make(map[uint]map[chan FileChangeEvent]struct{}),
	}
}

func SubscribeFileEvents(projectID uint) (<-chan FileChangeEvent, func()) {
	ch := make(chan FileChangeEvent, 8)

	globalFileEventBroker.mu.Lock()
	if globalFileEventBroker.subscribers[projectID] == nil {
		globalFileEventBroker.subscribers[projectID] = make(map[chan FileChangeEvent]struct{})
	}
	globalFileEventBroker.subscribers[projectID][ch] = struct{}{}
	globalFileEventBroker.mu.Unlock()

	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			globalFileEventBroker.mu.Lock()
			defer globalFileEventBroker.mu.Unlock()

			subscribers := globalFileEventBroker.subscribers[projectID]
			if subscribers == nil {
				return
			}

			delete(subscribers, ch)
			if len(subscribers) == 0 {
				delete(globalFileEventBroker.subscribers, projectID)
			}
		})
	}

	return ch, unsubscribe
}

func PublishFileEvent(event FileChangeEvent) {
	globalFileEventBroker.mu.RLock()
	subscribers := globalFileEventBroker.subscribers[event.ProjectID]
	if len(subscribers) == 0 {
		globalFileEventBroker.mu.RUnlock()
		return
	}

	targets := make([]chan FileChangeEvent, 0, len(subscribers))
	for subscriber := range subscribers {
		targets = append(targets, subscriber)
	}
	globalFileEventBroker.mu.RUnlock()

	for _, subscriber := range targets {
		select {
		case subscriber <- event:
		default:
		}
	}
}