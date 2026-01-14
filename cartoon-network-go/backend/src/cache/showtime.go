package cache

import (
	"sync"
	"time"

	"backend/db"
	"backend/models"
)

var (
	showTimeCache []models.Cartoon
	showMutex     sync.RWMutex
)

func RefreshShowTimeCache() {
	for {
		var cartoons []models.Cartoon
		db.DB.Select("id, name, show_time").Find(&cartoons)

		showMutex.Lock()
		showTimeCache = cartoons
		showMutex.Unlock()

		time.Sleep(30 * time.Second)
	}
}

func GetShowTimeCache() []models.Cartoon {
	showMutex.RLock()
	defer showMutex.RUnlock()
	return showTimeCache
}
