package worker

import (
	"log"

	"cartoon-network-go/backend/src/db"
	"cartoon-network-go/backend/src/models"
)

type Job struct {
	Type string
	Data interface{}
}

var JobQueue = make(chan Job, 200)

func StartWorkerPool() {
	for i := 1; i <= 6; i++ {
		go worker()
	}
}

func worker() {
	for job := range JobQueue {
		switch job.Type {

		case "LIKE":
			like := job.Data.(*models.Like)

			var count int64
			db.DB.Model(&models.Like{}).
				Where("user_id = ? AND cartoon_id = ?", like.UserID, like.CartoonID).
				Count(&count)

			if count == 0 {
				if err := db.DB.Create(like).Error; err != nil {
					log.Println("Duplicate like prevented")
				}
			}

		case "VIEW":
			db.DB.Create(job.Data.(*models.CartoonView))
		}
	}
}
