package worker

import (
	"backend/db"
	"backend/models"
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

		//  Views are async
		case "VIEW":
			db.DB.Create(job.Data.(*models.CartoonView))

		}
	}
}
