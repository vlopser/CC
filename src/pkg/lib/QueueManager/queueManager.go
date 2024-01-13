package queuemanager

import (
	"cc/pkg/models/task"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

const (
	REQUEST_QUEUE = "request_queue"
	WORKERS_GROUP = "workers_group"
)

func EnqueueTask(task task.Task, nats_server *nats.Conn) error {

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	err = nats_server.Publish(REQUEST_QUEUE, taskJSON)
	if err != nil {
		return err
	}

	return nil
}

func SubscribeQueueTask(nats_server *nats.Conn, callback func(task.Task, *nats.Conn)) {
	nats_server.QueueSubscribe(
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
}