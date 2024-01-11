package taskmanager

import (
	. "cc/pkg/lib/QueueManager"
	. "cc/pkg/lib/StoreManager"
	"cc/pkg/models/result"
	"cc/pkg/models/task"
	"log"
	"time"

	"github.com/google/uuid"
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

func PostResult(nats_server *nats.Conn, taskId uuid.UUID, output string, time_elapsed time.Duration) {
	result := result.Result{
		TaskId:      taskId,
		Output:      output, // Quitar?
		TimeElapsed: time_elapsed,
	}

	err := StoreResult(nats_server, result)
	if err != nil {
		log.Println("Error when storing the result of", taskId.String(), ":", err)
		return
	}
}

func GetTasks(nats_server *nats.Conn, handleFunc func(task.Task, *nats.Conn)) {
	SubscribeQueueTask(nats_server, handleFunc)
}
