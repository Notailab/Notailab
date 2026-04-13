package service

import (
	"container/list"
	"testing"
)

func newTestAgentService() *AgentService {
	return &AgentService{
		agentCache:      make(map[string]*list.Element),
		agentOrder:      list.New(),
		agentCacheLimit: 2,
	}
}

func TestAgentCacheMovesHitToFront(t *testing.T) {
	service := newTestAgentService()
	service.putCachedAgentLocked("a", &cachedAgent{})
	service.putCachedAgentLocked("b", &cachedAgent{})

	if got := service.getCachedAgentLocked("a"); got == nil {
		t.Fatal("expected cache hit for a")
	}

	front := service.agentOrder.Front()
	if front == nil {
		t.Fatal("expected cache order to have front element")
	}
	entry, ok := front.Value.(*cachedAgentEntry)
	if !ok || entry == nil {
		t.Fatal("unexpected cache entry type")
	}
	if entry.key != "a" {
		t.Fatalf("expected a to be most recent, got %q", entry.key)
	}
}

func TestAgentCacheEvictsLeastRecentlyUsed(t *testing.T) {
	service := newTestAgentService()
	service.putCachedAgentLocked("a", &cachedAgent{})
	service.putCachedAgentLocked("b", &cachedAgent{})

	if got := service.getCachedAgentLocked("a"); got == nil {
		t.Fatal("expected cache hit for a")
	}

	service.putCachedAgentLocked("c", &cachedAgent{})

	if service.getCachedAgentLocked("b") != nil {
		t.Fatal("expected b to be evicted as least recently used")
	}
	if service.getCachedAgentLocked("a") == nil {
		t.Fatal("expected a to remain cached")
	}
	if service.getCachedAgentLocked("c") == nil {
		t.Fatal("expected c to remain cached")
	}
}
