package repeated_submit

import (
	"sync"
	"time"
)

var (
	lockMap        = make(map[string]time.Time)
	lockMapMutex   sync.Mutex
	lockExpiration = 10 * time.Second
)

type DefRepeatedSubmitLock struct {
}

func NewDefRepeatedSubmitLock() RepeatedSubmitLock {
	return &DefRepeatedSubmitLock{}
}

func (d DefRepeatedSubmitLock) AcquireLock(key string) bool {
	lockMapMutex.Lock()
	defer lockMapMutex.Unlock()

	now := time.Now()
	if expiration, exists := lockMap[key]; exists {
		// If the lock exists and is not expired
		if now.Before(expiration) {
			return false
		}
		// Remove expired lock
		delete(lockMap, key)
	}

	// Set a new lock with expiration time
	lockMap[key] = now.Add(lockExpiration)
	return true
}

func (d DefRepeatedSubmitLock) ReleaseLock(key string) {
	lockMapMutex.Lock()
	defer lockMapMutex.Unlock()

	delete(lockMap, key)
}
