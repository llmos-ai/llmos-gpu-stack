package gpudevice

import (
	"sync"
)

type NodeDeviceInfo struct {
	Annotations map[string]string
}

// NodeDeviceThreadSafeCache is a structure to hold the cache
type NodeDeviceThreadSafeCache struct {
	cache map[string]NodeDeviceInfo
	mutex sync.RWMutex
}

// NewThreadSafeCache initializes a new cache of NodeDeviceInfo
func NewThreadSafeCache() *NodeDeviceThreadSafeCache {
	return &NodeDeviceThreadSafeCache{
		cache: make(map[string]NodeDeviceInfo),
	}
}

// Get retrieves a value from the cache in a thread-safe way
func (c *NodeDeviceThreadSafeCache) Get(key string) (NodeDeviceInfo, bool) {
	c.mutex.RLock()         // Acquire a read lock
	defer c.mutex.RUnlock() // Release the read lock

	value, exists := c.cache[key]
	return value, exists
}

// Set adds or updates a value in the cache in a thread-safe way
func (c *NodeDeviceThreadSafeCache) Set(key string, value NodeDeviceInfo) {
	c.mutex.Lock()         // Acquire a write lock
	defer c.mutex.Unlock() // Release the write lock

	c.cache[key] = value
}

// Delete removes a key from the cache in a thread-safe way
func (c *NodeDeviceThreadSafeCache) Delete(key string) {
	c.mutex.Lock()         // Acquire a write lock
	defer c.mutex.Unlock() // Release the write lock

	delete(c.cache, key)
}
