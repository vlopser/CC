package taskmanager

import (
	. "cc/src/pkg/lib/QueueManager"
	store "cc/src/pkg/lib/StoreManager"
	"cc/src/pkg/models/result"
	"cc/src/pkg/models/task"
	"log"
	"path"

	"github.com/nats-io/nats.go"

	"cc/src/pkg/models/request"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	TASK_ID_PARAM = "taskId"
)

func SetTaskStatusToExecuting(nats_server *nats.Conn, t task.Task) error {
	err := store.SetTaskStatus(nats_server, t.UserId, t.TaskId.String(), task.EXECUTING)
	if err != nil {
		log.Println("Error when changing the state of", t.TaskId.String(), "to executing:", err)
		return err
	}

	return nil
}

func SetTaskStatusToFinishedWithErrors(nats_server *nats.Conn, t task.Task) error {
	err := store.SetTaskStatus(nats_server, t.UserId, t.TaskId.String(), task.FINISHED_ERRORS)
	if err != nil {
		log.Println("Error when changing the state of", t.TaskId.String(), "to executing:", err)
		return err
	}

	return nil
}

func SetTaskStatusToFinished(nats_server *nats.Conn, t task.Task) error {
	err := store.SetTaskStatus(nats_server, t.UserId, t.TaskId.String(), task.FINISHED)
	if err != nil {
		log.Println("Error when changing the state of", t.TaskId.String(), "to finished:", err)
		return err
	}

	return nil
}

func SetTaskStatusToPending(nats_server *nats.Conn, t task.Task) error {
	err := store.SetTaskStatus(nats_server, t.UserId, t.TaskId.String(), task.PENDING)
	if err != nil {
		log.Println("Error when changing the state of", t.TaskId.String(), "to pending:", err)
		return err
	}

	return nil
}

func PostResult(nats_server *nats.Conn, result result.Result) {

	bucket := result.TaskId.String()
	err := store.CreateTaskBucket(nats_server, bucket)
	if err != nil {
		log.Println("Error when creating the bucket", result.TaskId.String(), ":", err)
		return
	}

	for _, file := range result.Files {
		err = store.StoreFileInBucket(nats_server, file, path.Base(file), bucket)
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

	log.Println(context.Request.Header)

	task := task.Task{
		TaskId:     uuid.New(),
		UserId:     context.Request.Header.Get("X-Forwarded-User"),
		RepoUrl:    requestBody.Url,
		Parameters: requestBody.Parameters,
	}

	err := SetTaskStatusToPending(nats_server, task)
	if err != nil {
		log.Println(err.Error())
		context.JSON(http.StatusInternalServerError, "An internal error happened.")
		return
	}

	err = EnqueueTask(task, nats_server)
	if err != nil {
		log.Println("Error enqueueing task", task.TaskId, ":", err.Error())
		context.JSON(http.StatusInternalServerError, "An internal error happened.")
		return
	}

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

	store.GetResult(nats_server, taskId)

	//zip con los res.File

	//context.FileAttachment(filePath, zip)

	// filePath := "hola.txt"
	// context.FileAttachment(filePath, "hola.txt")

	context.IndentedJSON(http.StatusOK, "ok")
}

func GetTaskStatus(context *gin.Context, nats_server *nats.Conn) {

	taskId := context.Request.URL.Query().Get(TASK_ID_PARAM)
	if taskId == "" {
		log.Println("Error: parameter taskId is missing")
		context.IndentedJSON(http.StatusBadRequest, "Error: parameter taskId is missing")
		return
	}

	userId := context.Request.Header.Get("X-Forwarded-User")

	task_state, err := store.GetTaskStatus(nats_server, taskId, userId)
	if err != nil {
		context.IndentedJSON(http.StatusForbidden, "Error: no task with given ID")
	}

	context.IndentedJSON(http.StatusOK, task_state.String())
}

func GetAllTasks(context *gin.Context, nats_server *nats.Conn) {

	userId := context.Request.Header.Get("X-Forwarded-User")

	allTaskIds, _ := store.GetUserTasks(nats_server, userId)
	// if err != nil {
	// 	context.IndentedJSON(http.StatusForbidden, err.Error())
	// }

	if len(allTaskIds) == 0 {
		context.IndentedJSON(http.StatusNoContent, allTaskIds)
	}

	context.IndentedJSON(http.StatusOK, allTaskIds)
}
