package taskmanager

import (
	"cc/pkg/models/request"
	"cc/pkg/models/task"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

const (
	TASK_ID_PARAM = "taskId"
)

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
