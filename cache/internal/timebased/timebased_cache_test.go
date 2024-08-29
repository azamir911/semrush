package timebased

import (
	"testing"
	"time"
)

func TestTimeBasedCache_SetAndGet(t *testing.T) {
	cache := NewTimeBasedCache(2, 2*time.Second)

	cache.Set("a", 1)
	cache.Set("b", 2)

	if val, exists := cache.Get("a"); !exists || val != 1 {
		t.Errorf("Expected key 'a' to exist with value 1, got %v", val)
	}

	time.Sleep(3 * time.Second)

	if _, exists := cache.Get("a"); exists {
		t.Errorf("Expected key 'a' to be expired and not exist")
	}
}

func TestTimeBasedCache_Eviction(t *testing.T) {
	cache := NewTimeBasedCache(2, 5*time.Second)

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3) // This should evict "a"

	if _, exists := cache.Get("a"); exists {
		t.Errorf("Expected key 'a' to be evicted, but it still exists")
	}

	if val, exists := cache.Get("b"); !exists || val != 2 {
		t.Errorf("Expected key 'b' to exist with value 2, got %v", val)
	}
}

func TestTimeBasedCache_Cleanup(t *testing.T) {
	cache := NewTimeBasedCache(2, 2*time.Second)

	cache.Set("a", 1)
	cache.Set("b", 2)

	time.Sleep(3 * time.Second)

	cache.Mutex.RLock()
	if len(cache.ttlMap) != 0 {
		t.Errorf("Expected cache to be cleaned up, but it still has items")
	}
	cache.Mutex.RUnlock()
}

func TestTimeBasedCache_StopCleanup(t *testing.T) {
	cache := NewTimeBasedCache(2, 2*time.Second)
	defer cache.StopCleanup() // Ensure cleanup goroutine is stopped when test ends

	cache.Set("a", 1)
	cache.Set("b", 2)

	// Allow time for some items to be cleaned up
	time.Sleep(1 * time.Second)

	// Manually check the cache's internal state
	cache.Mutex.RLock()
	if len(cache.ttlMap) != 2 {
		t.Errorf("Expected cache to have 2 items, but it has %d", len(cache.ttlMap))
	}
	cache.Mutex.RUnlock()

	// Stop the cleanup process
	cache.StopCleanup()

	// Wait for a duration longer than the expiration to ensure no more items are cleaned up
	time.Sleep(3 * time.Second)

	// Check cache state after stopping cleanup
	cache.Mutex.RLock()
	if len(cache.ttlMap) == 0 {
		t.Errorf("Expected cache to retain items after stopping cleanup, but it has %d", len(cache.ttlMap))
	}
	cache.Mutex.RUnlock()
}

func TestTimeBasedCache_EvictionWithStop(t *testing.T) {
	cache := NewTimeBasedCache(2, 2*time.Second)
	defer cache.StopCleanup() // Ensure cleanup goroutine is stopped when test ends

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3) // This should evict "a"

	time.Sleep(1 * time.Second) // Wait for eviction

	cache.Mutex.RLock()
	if _, exists := cache.Data["a"]; exists {
		t.Errorf("Expected key 'a' to be evicted, but it still exists")
	}
	cache.Mutex.RUnlock()

	// Stop the cleanup process
	cache.StopCleanup()

	// Check cache state after stopping cleanup
	cache.Mutex.RLock()
	if len(cache.Data) != 2 {
		t.Errorf("Expected cache to retain 2 items after stopping cleanup, but it has %d", len(cache.Data))
	}
	cache.Mutex.RUnlock()
}

func TestTimeBasedCache_ImmediateExpiration(t *testing.T) {
	cache := NewTimeBasedCache(2, 1*time.Second)
	defer cache.StopCleanup() // Ensure cleanup goroutine is stopped when test ends

	cache.Set("a", 1)

	// Wait for the item to expire
	time.Sleep(2 * time.Second)

	cache.Mutex.RLock()
	if _, exists := cache.Data["a"]; exists {
		t.Errorf("Expected key 'a' to be expired, but it still exists")
	}
	cache.Mutex.RUnlock()
}
