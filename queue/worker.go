package queue

import (
	"context"
	"sync"
	"time"

	db "github.com/ssubedir/go-shoot/db"

	"go.mongodb.org/mongo-driver/bson"
)

// Worker - Worker that procresses tasks
type Worker struct {

	// Task Channals
	ReadyChan    chan chan Task
	AssignedTask chan Task

	// Worker synchronization
	IsDone *sync.WaitGroup

	// Worker quit
	Quit chan bool

	DB *db.MongoClient
}

// NewWorker - Creates a new worker
func NewWorker(readyPool chan chan Task, done *sync.WaitGroup) *Worker {

	// DB
	Connections := db.DBConnection()
	mongo := Connections["mongodb"].(db.Mongo)
	mongoConn := mongo.Connect()

	// Return New Worker
	return &Worker{
		ReadyChan:    readyPool,
		AssignedTask: make(chan Task),
		IsDone:       done,
		Quit:         make(chan bool),
		DB:           mongoConn,
	}
}

// Start - Begins processing worker's task
func (w *Worker) Start() {
	go func() {
		w.IsDone.Add(1)
		for {
			w.ReadyChan <- w.AssignedTask // check the Task queue in
			select {
			case Task := <-w.AssignedTask: // see if anything has been assigned to the queue
				// Task enqueued
				w.enqueueState(Task)

				start := time.Now()
				Task.Run()
				elapsed := time.Since(start)
				Task.metadata.Duration = elapsed
				w.succeededState(Task)

				// if recover() != nil {
				// 	// error run task again
				// 	Task.Run()
				// }

			case <-w.Quit:
				w.DB.Disconnect()
				w.IsDone.Done()
				return
			}
		}
	}()
}

// Stop - Stops the worker
func (w *Worker) Stop() {
	w.Quit <- true
}

func (w *Worker) enqueueState(task Task) {

	// update Db
	client := w.DB.MONGO
	goshoot := client.Database("go-shoot")
	enqueued := goshoot.Collection("enqueued")
	enqueued.InsertOne(context.TODO(),
		bson.M{
			"task_id":            task.id,
			"task_status":        "enqueued",
			"task_name":          task.metadata.Name,
			"task_type":          task.metadata.TaskType,
			"task_enqueued_time": time.Now(),
		})
}

func (w *Worker) succeededState(task Task) {

	// update Db
	client := w.DB.MONGO
	goshoot := client.Database("go-shoot")
	enqueued := goshoot.Collection("enqueued")
	succeeded := goshoot.Collection("succeeded")

	// Queue update
	enqueued.DeleteOne(context.TODO(), bson.M{"task_id": task.id})

	// Scheduled Update
	if task.metadata.TaskType == "Scheduled" {
		scheduled := goshoot.Collection("scheduled")
		scheduled.DeleteOne(context.TODO(), bson.M{"task_id": task.id})
	}

	// Success
	succeeded.InsertOne(context.TODO(),
		bson.M{
			"task_id":             task.id,
			"task_status":         "succeeded",
			"task_type":           task.metadata.TaskType,
			"task_duration":       task.metadata.Duration,
			"task_name":           task.metadata.Name,
			"task_succeeded_time": time.Now(),
		})

}
