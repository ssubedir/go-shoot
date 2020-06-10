# About Go Shoot
An easy way to perform background processing in go. Open and free.  [Go Shoot Dashboard](https://github.com/ssubedir/go-shoot-dashboard)

## Getting Started

These instructions will get you a copy of go-shoot and instruct on integrating into your projects.

### Installation

Install go-shoot using ```go get```

```sh
$ go get github.com/ssubedir/go-shoot
```
# Setup
Create a .env file with the following Environmental Variables

```
# MongoDB info

DB_USER=username

DB_PASSWORD=password

DB_CLUSTER=address/cluster
```
# Example
Import this package and write
```go
package main
import (
	"fmt"
	q "github.com/ssubedir/go-shoot/queue"
)

func  TestTask() {
	fmt.Println("Task c")
}
func main() {
	
	var  a = &q.Task{Run: func() {
		fmt.Println("Task a ")
		time.Sleep(5 * time.Second) // heavy task
	}}
	var  b = &q.Task{Run: func() {
		fmt.Println("Task b ")
	}}

	var  c = &q.Task{Run: TestTask} // task c
	
	tasks := q.NewQueue(8)
	tasks.Start()
	defer tasks.Stop()

	tasks.Enqueue(&q.BackgroundTask{Task: a, Name: "Task a"})
	tasks.Enqueue(&q.BackgroundTask{Task: b, Name: "Task b"})
	tasks.Schedule(&q.ScheduledTask{Task: c, Name: "Task c"}, "5s") // delayed task

	...
	...
	...
}

```
  [More Examples](https://github.com/ssubedir/go-shoot)

## License

This project is licensed under the MIT License - see the [LICENSE.md](https://github.com/ssubedir/go-shoot/blob/master/LICENSE) file for details


