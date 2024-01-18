package taskmanager

import (
	. "cc/src/pkg/lib/QueueManager"
	store "cc/src/pkg/lib/StoreManager"
	"cc/src/pkg/models/errors"
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
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if requestBody.Parameters == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'parameters' is missing."})
		return
	}
	if requestBody.Url == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'url' is missing."})
		return
	}

	task := task.Task{
		TaskId:     uuid.New(),
		UserId:     context.Request.Header.Get("X-Forwarded-User"),
		RepoUrl:    requestBody.Url,
		Parameters: requestBody.Parameters,
	}

	err := SetTaskStatusToPending(nats_server, task)
	switch err {
	case nil:
		break
	case errors.ErrUserInvalid:
		context.JSON(http.StatusUnauthorized, gin.H{"error": "User ID is invalid."})
		return
	case errors.ErrTaskInvalid:
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Task ID is invalid."})
		return
	default:
		context.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error happened."})
		return
	}

	err = EnqueueTask(task, nats_server)
	if err != nil {
		context.JSON(http.StatusInternalServerError, "An internal error happened.")
		return
	}

	// todo llamar la libreria
	context.JSON(http.StatusCreated, gin.H{"data": task.TaskId})
}

func GetTaskResult(context *gin.Context, nats_server *nats.Conn) {
	queryParams := context.Request.URL.Query()
	taskId := queryParams.Get(TASK_ID_PARAM)
	if taskId == "" {
		log.Println("Error: parameter taskId is missing")
		context.JSON(http.StatusBadRequest, gin.H{"error": "Parameter taskId is missing"})
		return
	}

	log.Println("Received request to get result for task " + taskId)

	_, err := store.GetResult(nats_server, taskId)
	switch err {
	case nil:
		break

	case errors.ErrTaskInvalid:
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Task ID is invalid."})
		return

	case errors.ErrTaskNotFound:
		context.JSON(http.StatusBadRequest, gin.H{"error": "Given task ID does not exist."})
		return

	default:
		context.JSON(http.StatusInternalServerError, "An internal error happened.")
	}

	//zip con los res.File

	//context.FileAttachment(filePath, zip)

	// filePath := "hola.txt"
	// context.FileAttachment(filePath, "hola.txt")

	context.IndentedJSON(http.StatusOK, "ok")
}

func GetTaskStatus(context *gin.Context, nats_server *nats.Conn) {

	taskId := context.Request.URL.Query().Get(TASK_ID_PARAM)
	if taskId == "" {
		context.JSON(http.StatusBadRequest, "Error: parameter taskId is missing")
		return
	}

	userId := context.Request.Header.Get("X-Forwarded-User")

	task_state, err := store.GetTaskStatus(nats_server, taskId, userId)
	switch err {
	case nil:
		break

	case errors.ErrUserNotFound, errors.ErrTaskNotFound:
		context.JSON(http.StatusBadRequest, gin.H{"error": "Given task does not exist."})
		return

	case errors.ErrUserInvalid:
		context.JSON(http.StatusUnauthorized, gin.H{"error": "User ID is invalid."})
		return

	default:
		context.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error has happened."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": task_state.String()})
}

func GetAllTasks(context *gin.Context, nats_server *nats.Conn) {

	userId := context.Request.Header.Get("X-Forwarded-User")

	allTaskIds, err := store.GetUserTasks(nats_server, userId)
	switch err {
	case nil:
		break

	case errors.ErrUserNotFound:
		context.JSON(http.StatusOK, gin.H{"data": []string{}})
		return

	case errors.ErrUserInvalid:
		context.JSON(http.StatusUnauthorized, gin.H{"error": "User ID is invalid."})
		return

	default:
		context.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error has happened."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": allTaskIds})
}
