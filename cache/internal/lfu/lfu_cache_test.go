package lfu

import "testing"

func TestLFUCache_EvictionOrder(t *testing.T) {
	cache := NewLFUCache(2)

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Get("a")    // Access 'a', so its frequency increases
	cache.Set("c", 3) // This should evict 'b', as it has the lowest frequency

	if _, exists := cache.Get("b"); exists {
		t.Errorf("Expected key 'b' to be evicted, but it still exists.")
	}

	if _, exists := cache.Get("a"); !exists {
		t.Errorf("Expected key 'a' to exist, but it was evicted.")
	}

	if _, exists := cache.Get("c"); !exists {
		t.Errorf("Expected key 'c' to exist, but it was evicted.")
	}
}

func TestLFUCache_UpdateExistingKey(t *testing.T) {
	cache := NewLFUCache(2)

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("a", 3) // Update 'a' and reset its frequency to 1

	value, _ := cache.Get("a")
	if value != 3 {
		t.Errorf("Expected value of key 'a' to be 3, got %v", value)
	}

	cache.Set("c", 4) // This should evict 'b', since 'a' was recently updated and accessed

	if _, exists := cache.Get("b"); exists {
		t.Errorf("Expected key 'b' to be evicted, but it still exists.")
	}
}

func TestLFUCache_ExistenceCheck(t *testing.T) {
	cache := NewLFUCache(3)

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)

	if _, exists := cache.Get("a"); !exists {
		t.Errorf("Expected key 'a' to exist, but it does not.")
	}

	if _, exists := cache.Get("b"); !exists {
		t.Errorf("Expected key 'b' to exist, but it does not.")
	}

	if _, exists := cache.Get("c"); !exists {
		t.Errorf("Expected key 'c' to exist, but it does not.")
	}

	cache.Set("d", 4) // This should evict 'a' since it's the least frequently used

	if _, exists := cache.Get("a"); exists {
		t.Errorf("Expected key 'a' to be evicted, but it still exists.")
	}
}
