package taskmanager

import (
	models "cc/Models"
	. "cc/QueueManager"
	. "cc/StoreManager"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

func SetTaskStatusToExecuting(nats_server *nats.Conn, taskId string) {
	err := ChangeState(nats_server, taskId, models.EXECUTING)
	if err != nil {
		log.Println("Error when changing the state of", taskId, "to executing:", err)
		return
	}
}

func SetTaskStatusToFinished(nats_server *nats.Conn, taskId string) {
	err := ChangeState(nats_server, taskId, models.FINISHED)
	if err != nil {
		log.Println("Error when changing the state of", taskId, "to finished:", err)
		return
	}
}

func SetTaskStatusToPending(nats_server *nats.Conn, taskId string) {
	err := ChangeState(nats_server, taskId, models.PENDING)
	if err != nil {
		log.Println("Error when changing the state of", taskId, "to pending:", err)
		return
	}
}

func PostResult(nats_server *nats.Conn, taskId uuid.UUID, output string) {
	result := models.Result{
		TaskId:    taskId,
		Output:    output, //cambiar a resultado
		Timestamp: time.Now(),
	}
	err := StoreResult(nats_server, result)
	if err != nil {
		log.Println("Error when storing the result of", taskId.String(), ":", err)
		return
	}
}

func GetTasks(nats_server *nats.Conn, handleFunc func(models.Task, *nats.Conn)) {
	SubscribeQueueTask(nats_server, handleFunc)
}
