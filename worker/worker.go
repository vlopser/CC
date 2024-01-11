package main

import (
	. "cc/pkg/lib/QueueManager"
	. "cc/pkg/lib/StoreManager"
	. "cc/pkg/lib/TaskManager"
	"cc/pkg/models/task"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

var wg sync.WaitGroup

func executeTask(task task.Task, nats_server *nats.Conn) {
	log.Println("Request received!")

	SetTaskStatusToExecuting(nats_server, task.TaskId.String())

	fmt.Println("Executed task") // result = func()

	PostResult(nats_server, task.TaskId, task.Input)

	SetTaskStatusToFinished(nats_server, task.TaskId.String())

	//CODIGO DEL FRONTEND
	// res, err := GetResult(nats_server, task.TaskId.String())
	// if err != nil {
	// 	log.Println("Error when storing the result of", task.TaskId, ":", err)
	// 	return
	// }
	// println(res.Output)
}

func waitForSigkill() {

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGINT)

	<-c

	wg.Done()

	os.Exit(0)
}

func waitForTasks(nats_server *nats.Conn) {
	GetTasks(nats_server, executeTask)

	wg = sync.WaitGroup{}
	wg.Add(1)

	log.Println("Waiting for request. (Presiona Ctrl+C para salir)")

	wg.Wait()                   // Esperar a que se reciba una señal de interrupción
	time.Sleep(1 * time.Second) // Permitir tiempo para desuscribirse antes de salir
}

func main() {

	nats_server, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nats_server.Close()

	go waitForSigkill()

	// GetTasks(nats_server, executeTask)

	// task := task.Task{
	// 	TaskId: uuid.New(),
	// 	Input:  "www.upv.es",
	// 	Status: task.PENDING,
	// }

	// EnqueueTask(task, nats_server)

	waitForTasks(nats_server)
}
