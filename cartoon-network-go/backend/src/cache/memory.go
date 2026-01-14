package cache

import (
	"sync"
	"time"

	"backend/db"
	"backend/models"
)

var (
	homeCache []models.Cartoon
	mutex     sync.RWMutex
)

// RefreshHomeCache continuously refreshes home page cartoon cache
func RefreshHomeCache() {
	for {
		var cartoons []models.Cartoon

		db.DB.
			Preload("Images").
			Preload("Characters").
			Find(&cartoons)

		mutex.Lock()
		homeCache = cartoons
		mutex.Unlock()

		time.Sleep(30 * time.Second)
	}
}

// GetHomeCache safely returns cached cartoons
func GetHomeCache() []models.Cartoon {
	mutex.RLock()
	defer mutex.RUnlock()
	return homeCache
}
