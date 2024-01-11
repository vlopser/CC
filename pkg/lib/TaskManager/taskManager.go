package taskmanager

import (
	. "cc/pkg/lib/QueueManager"
	. "cc/pkg/lib/StoreManager"
	"cc/pkg/models/result"
	"cc/pkg/models/task"
	"log"

	"github.com/nats-io/nats.go"
)

func SetTaskStatusToExecuting(nats_server *nats.Conn, taskId string) {
	err := ChangeState(nats_server, taskId, task.EXECUTING)
	if err != nil {
		log.Println("Error when changing the state of", taskId, "to executing:", err)
		return
	}
}

func SetTaskStatusToFinished(nats_server *nats.Conn, taskId string) {
	err := ChangeState(nats_server, taskId, task.FINISHED)
	if err != nil {
		log.Println("Error when changing the state of", taskId, "to finished:", err)
		return
	}
}

func SetTaskStatusToPending(nats_server *nats.Conn, taskId string) {
	err := ChangeState(nats_server, taskId, task.PENDING)
	if err != nil {
		log.Println("Error when changing the state of", taskId, "to pending:", err)
		return
	}
}

func PostResult(nats_server *nats.Conn, result result.Result) {

	bucket := result.TaskId.String()
	err := CreateTaskBucket(nats_server, bucket)
	if err != nil {
		log.Println("Error when creating the bucket", result.TaskId.String(), ":", err)
		return
	}

	file_name := "output.txt"
	file_path := result.TaskId.String() + "/result/" + file_name
	err = StoreFileInBucket(nats_server, file_path, file_name, bucket)
	if err != nil {
		log.Println("Error when storing the result of", result.TaskId.String(), ":", err)
		return
	}

	file_name = "errors.txt"
	file_path = result.TaskId.String() + "/result/" + file_name
	err = StoreFileInBucket(nats_server, file_path, file_name, bucket)
	if err != nil {
		log.Println("Error when storing the result of", result.TaskId.String(), ":", err)
		return
	}
}

func GetTasks(nats_server *nats.Conn, handleFunc func(task.Task, *nats.Conn)) {
	SubscribeQueueTask(nats_server, handleFunc)
}
