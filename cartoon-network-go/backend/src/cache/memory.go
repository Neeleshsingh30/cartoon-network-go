package cache

import (
	"sync"
	"time"

	"cartoon-network-go/backend/src/db"
	"cartoon-network-go/backend/src/models"
)

var (
	homeCache []models.Cartoon
	mutex     sync.RWMutex
)

func RefreshHomeCache() {
	for {
		var cartoons []models.Cartoon
		db.DB.Preload("Images").Find(&cartoons)

		mutex.Lock()
		homeCache = cartoons
		mutex.Unlock()

		time.Sleep(30 * time.Second)
	}
}

func GetHomeCache() []models.Cartoon {
	mutex.RLock()
	defer mutex.RUnlock()
	return homeCache
}
