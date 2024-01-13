package taskmanager

import (
	. "cc/pkg/lib/QueueManager"
	. "cc/pkg/lib/StoreManager"
	"cc/pkg/models/result"
	"cc/pkg/models/task"
	"log"
	"path"

	"github.com/nats-io/nats.go"
  
  "cc/pkg/models/request"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	TASK_ID_PARAM = "taskId"
)

func SetTaskStatusToExecuting(nats_server *nats.Conn, taskId string) {
	err := ChangeState(nats_server, taskId, task.EXECUTING)
	if err != nil {
		log.Println("Error when changing the state of", taskId, "to executing:", err)
		return
	}
}

func SetTaskStatusToFinishedWithErrors(nats_server *nats.Conn, taskId string) {
	err := ChangeState(nats_server, taskId, task.FINISHED_ERRORS)
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

	for _, file := range result.Files {
		err = StoreFileInBucket(nats_server, file, path.Base(file), bucket)
		// if err != nil {
		// 	log.Println("Error when storing the file", file, ":", err)
		// 	return
		// }
	}
}

func GetTasks(nats_server *nats.Conn, handleFunc func(task.Task, *nats.Conn)) {
	SubscribeQueueTask(nats_server, handleFunc)
}

func CreateTask(context *gin.Context, nats_server *nats.Conn) {

	var requestBody request.PostTaskBody
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Received request to create a task from " + requestBody.Url)
	log.Println("Task parameters are:")
	for _, param := range requestBody.Parameters {
		log.Println(param)
	}

	task := task.Task{
		TaskId:     uuid.New(),
		UserMail:   context.Request.Header.Get("X-Forwarded-Email"),
		RepoUrl:    requestBody.Url,
		Parameters: requestBody.Parameters,
		Status:     100,
	}

	//EnqueueTask(task, nats_server)

	// todo llamar la libreria
	context.IndentedJSON(http.StatusCreated, task.TaskId)
}

func GetTaskResult(context *gin.Context, nats_server *nats.Conn) {
	queryParams := context.Request.URL.Query()
	taskId := queryParams.Get(TASK_ID_PARAM)
	if taskId == "" {
		log.Println("Error: parameter taskId is missing")
		context.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	log.Println("Received request to get result for task " + taskId)

	//res := GetResult(nats_server, taskId)

	//zip con los res.File

	//context.FileAttachment(filePath, zip)

	// filePath := "hola.txt"
	// context.FileAttachment(filePath, "hola.txt")

	context.IndentedJSON(http.StatusCreated, "ok")
}