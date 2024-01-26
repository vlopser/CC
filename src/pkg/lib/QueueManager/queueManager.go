package queuemanager

import (
	"cc/src/pkg/models/requestInjection"
	"cc/src/pkg/models/task"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

const (
	TASK_QUEUE    = "task_queue"
	WORKERS_GROUP = "workers_group"

	INJECTION_QUEUE = "injection_queue"
	INJECTORS_GROUP = "injectors_group"
)

func EnqueueTask(task task.Task, nats_server *nats.Conn) error {

	taskJSON, err := json.Marshal(task)
	if err != nil {
		log.Println("Error marshalling task", task.TaskId, ":", err.Error())
		return err
	}

	err = nats_server.Publish(TASK_QUEUE, taskJSON)
	if err != nil {
		log.Println("Error enqueueing task", task.TaskId, ":", err.Error())
		return err
	}

	return nil
}

func EnqueueInjectionRequest(req requestInjection.RequestInjection, nats_server *nats.Conn) error {

	taskJSON, err := json.Marshal(req)
	if err != nil {
		log.Println("Error marshalling injection request:", err.Error())
		return err
	}

	err = nats_server.Publish(INJECTION_QUEUE, taskJSON)
	if err != nil {
		log.Println("Error enqueueing injection request:", err.Error())
		return err
	}

	return nil
}

func SubscribeQueueTask(nats_server *nats.Conn, callback func(task.Task, *nats.Conn)) error {
	_, err := nats_server.QueueSubscribe(
		TASK_QUEUE,
		WORKERS_GROUP,
		func(msg *nats.Msg) {
			fmt.Printf("Mensaje recibido: %s\n", msg.Data)

			var task task.Task
			err := json.Unmarshal(msg.Data, &task)
			if err != nil {
				fmt.Println("Error when unmarshalling JSON:", err)
				return
			}

			callback(task, nats_server)

		})
	if err != nil {
		log.Fatal("Error when subscribing task queue:", err.Error())
	}

	return nil
}

func SubscribeQueueInjection(nats_server *nats.Conn, callback func(*nats.Conn, requestInjection.RequestInjection)) error {
	_, err := nats_server.QueueSubscribe(
		INJECTION_QUEUE,
		INJECTORS_GROUP,
		func(msg *nats.Msg) {
			log.Println("Received request")
			var req requestInjection.RequestInjection
			err := json.Unmarshal(msg.Data, &req)
			if err != nil {
				fmt.Println("Error when unmarshalling JSON:", err)
				return
			}

			callback(nats_server, req)
		},
	)
	if err != nil {
		log.Fatal("Error when subscribing injector queue:", err.Error())
	}

	return nil
}
