package queuemanager

import (
	"cc/src/pkg/models/task"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

const (
	REQUEST_QUEUE = "request_queue"
	WORKERS_GROUP = "workers_group"
)

func EnqueueTask(task task.Task, nats_server *nats.Conn) error {

	taskJSON, err := json.Marshal(task)
	if err != nil {
		log.Println("Error marshalling task", task.TaskId, ":", err.Error())
		return err
	}

	err = nats_server.Publish(REQUEST_QUEUE, taskJSON)
	if err != nil {
		log.Println("Error enqueueing task", task.TaskId, ":", err.Error())
		return err
	}

	return nil
}

func SubscribeQueueTask(nats_server *nats.Conn, callback func(task.Task, *nats.Conn)) error {
	_, err := nats_server.QueueSubscribe(
		REQUEST_QUEUE,
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
		log.Println("Error when subscribing queue:", err.Error())
		return err
	}

	return nil

}
