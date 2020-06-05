package queue

import (
	"context"
	"sync"
	"time"

	db "github.com/ssubedir/go-shoot/db"

	"go.mongodb.org/mongo-driver/bson"
)

// Queue - A queue for enqueueing tasks to be processed
type Queue struct {
	// Channals
	QueueChan chan Task      // Task Channal
	ReadyChan chan chan Task // Ready Task Channals

	// Goroutine synchronization
	DispatcherSync *sync.WaitGroup // Work Dispatcher synchronization
	WorkersSync    *sync.WaitGroup // Workers synchronization

	// Queue Workers
	Workers []*Worker

	// Quit Queue
	Quit chan bool

	DB *db.MongoClient

	Name string
}

// NewQueue - Creates a new Queue
func NewQueue(n int) *Queue {

	// DB
	Connections := db.DBConnection()
	mongo := Connections["mongodb"].(db.Mongo)
	mongoConn := mongo.Connect()

	// workers
	w := make([]*Worker, n, n)

	// workers synchronization
	ws := sync.WaitGroup{}

	// Ready Task Channals
	rc := make(chan chan Task, n)

	// create n Workers
	for i := 0; i < n; i++ {
		w[i] = NewWorker(rc, &ws)
	}

	// return Queue
	return &Queue{
		// Channals
		QueueChan: make(chan Task),
		ReadyChan: rc,

		// Queue Workers
		Workers: w,

		// Goroutine synchronization
		DispatcherSync: &sync.WaitGroup{},
		WorkersSync:    &ws,

		// Quit Queue
		Quit: make(chan bool),

		DB: mongoConn,
	}
}

// dispatch - Dispatch workers to process tasks
func (q *Queue) dispatch() {
	q.DispatcherSync.Add(1)
	for {
		select {
		case Task := <-q.QueueChan: // We got something in on our queue
			workerChannel := <-q.ReadyChan // Check out an available worker
			workerChannel <- Task          // Send the request to the channel
		case <-q.Quit:
			for i := 0; i < len(q.Workers); i++ {
				q.Workers[i].Stop()
			}
			q.WorkersSync.Wait()
			q.DispatcherSync.Done()
			return
		}
	}
}

// Start - Starts the worker and dispatcher go routines
func (q *Queue) Start() {
	for i := 0; i < len(q.Workers); i++ {
		q.Workers[i].Start() // start workers
	}
	go q.dispatch() // queue dispach
}

// Stop - Stopes Queue
func (q *Queue) Stop() {
	q.Quit <- true          // pass quit flag
	q.DispatcherSync.Wait() // wait
}

// Enqueue - Fire-and-forget task are executed only once.
func (q *Queue) Enqueue(bt *BackgroundTask) {
	bt.GenerateMetaData(0)
	q.QueueChan <- *bt.Task
}

// Schedule - Delayed task are executed only once too, but not immediately, after a certain time interval.
func (q *Queue) Schedule(bt *ScheduledTask, durationstring string) {
	dur, _ := time.ParseDuration(durationstring)
	bt.GenerateMetaData(1)
	finished := make(chan bool)
	go func() {
		// update Db
		client := q.DB.MONGO
		goshoot := client.Database("go-shoot")
		schedule := goshoot.Collection("scheduled")
		schedule.InsertOne(context.TODO(),
			bson.M{
				"task_id":             bt.Task.id,
				"task_name":           bt.Name,
				"task_status":         "scheduled",
				"task_created_time":   time.Now(),
				"task_scheduled_time": time.Now().Add(dur)})
		t := time.NewTicker(dur)
		defer t.Stop()
		<-t.C
		q.QueueChan <- *bt.Task

		finished <- true
	}()
	<-finished
}
