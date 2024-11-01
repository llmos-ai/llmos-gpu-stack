package gpudevice

import (
	"reflect"
	"sync"
	"testing"
)

// TestSetAndGet tests setting and getting values in the cache
func TestSetAndGet(t *testing.T) {
	cache := NewThreadSafeCache()
	key := "testKey"
	value := NodeDeviceInfo{Annotations: map[string]string{
		"foo":                   "bar",
		"nvidia.com/gpu.device": "0",
	}}

	// Set a value
	cache.Set(key, value)

	// Get the value back and check if it exists
	retrievedValue, found := cache.Get(key)
	if !found {
		t.Fatalf("Expected to find key '%s' in cache", key)
	}
	if !reflect.DeepEqual(retrievedValue, value) {
		t.Fatalf("Expected value '%v', got '%v'", value, retrievedValue)
	}
}

// TestGetNonExistentKey tests retrieving a non-existent key
func TestGetNonExistentKey(t *testing.T) {
	cache := NewThreadSafeCache()
	key := "nonExistentKey"

	// Try to get a value that doesn't exist
	_, found := cache.Get(key)
	if found {
		t.Fatalf("Did not expect to find key '%s' in cache", key)
	}
}

// TestDelete tests deleting a key from the cache
func TestDelete(t *testing.T) {
	cache := NewThreadSafeCache()
	key := "deleteKey"
	value := NodeDeviceInfo{Annotations: map[string]string{
		"foo":                   "bar",
		"nvidia.com/gpu.device": "0",
	}}

	// Set a value and then delete it
	cache.Set(key, value)
	cache.Delete(key)

	// Ensure the key no longer exists
	_, found := cache.Get(key)
	if found {
		t.Fatalf("Expected key '%s' to be deleted from cache", key)
	}
}

// TestConcurrentAccess tests concurrent access to the cache
func TestConcurrentAccess(t *testing.T) {
	cache := NewThreadSafeCache()
	var wg sync.WaitGroup
	key := "concurrentKey"
	value := NodeDeviceInfo{Annotations: map[string]string{
		"foo": "bar",
	}}

	// Set up a number of goroutines to access the cache concurrently
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Set(key, value)
			retrievedValue, _ := cache.Get(key)
			if !reflect.DeepEqual(retrievedValue, value) {
				t.Errorf("Expected value '%v', got '%v'", value, retrievedValue)
			}
		}()
	}
	wg.Wait()

	// Final check to ensure the value is still correct
	retrievedValue, found := cache.Get(key)
	if !found {
		t.Fatalf("Expected to find key '%s' in cache after concurrent access", key)
	}
	if !reflect.DeepEqual(retrievedValue, value) {
		t.Fatalf("Expected value '%v', got '%v'", value, retrievedValue)
	}
}
