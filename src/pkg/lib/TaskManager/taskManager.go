package taskmanager

import (
	"bytes"
	. "cc/src/pkg/lib/QueueManager"
	store "cc/src/pkg/lib/StoreManager"
	"cc/src/pkg/models/errors"
	"cc/src/pkg/models/result"
	"cc/src/pkg/models/task"
	"cc/src/pkg/utils"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"cc/src/pkg/models/request"
	"net/http"
	"net/url"

	"archive/zip"

	"github.com/gin-gonic/gin"
)

const (
	TASK_ID_PARAM = "taskId"
)

/************************ TASK STATUS ************************/

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

func SetTaskStatusToExecuting(nats_server *nats.Conn, t task.Task) error {
	err := store.SetTaskStatus(nats_server, t.UserId, t.TaskId.String(), task.EXECUTING)
	if err != nil {
		log.Println("Error when changing the state of", t.TaskId.String(), "to executing:", err)
		return err
	}

	return nil
}

func SetTaskStatusToUnexpectedError(nats_server *nats.Conn, t task.Task) error {
	err := store.SetTaskStatus(nats_server, t.UserId, t.TaskId.String(), task.UNEXPECTED_ERROR)
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

/************************ TASK RESULTS ************************/

func CreateTaskResult(nats_server *nats.Conn, result result.Result) error {

	bucket := result.TaskId.String()
	err := store.CreateTaskBucket(nats_server, bucket)
	if err != nil {
		log.Println("Error when creating the bucket", result.TaskId.String(), ":", err)
		return err
	}

	for _, file := range result.Files {
		err = store.StoreFileInBucket(nats_server, file, path.Base(file), bucket)
		// if err != nil {
		// 	log.Println("Error when storing the file", file, ":", err)
		// 	return
		// }
	}

	return nil
}

func GetTaskResult(context *gin.Context, nats_server *nats.Conn) {
	queryParams := context.Request.URL.Query()
	if !checkGetTaskResult(context, nats_server, &queryParams) {
		return
	}

	taskId := queryParams.Get(TASK_ID_PARAM)
	res, err := store.GetResult(nats_server, taskId)
	switch err {
	case nil:
		break

	case errors.ErrTaskInvalid:
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Task ID is invalid."})
		return

	case errors.ErrTaskNotFound:
		context.JSON(http.StatusBadRequest, gin.H{"error": "No results for given task ID"})
		return

	default:
		context.JSON(http.StatusInternalServerError, "An internal error happened.")
		return
	}
	defer utils.RemoveAllFiles(res.Files)

	//Create a zip with result files
	outputZip := res.TaskId.String() + "_results.zip"

	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)

	for _, file := range res.Files {
		err := utils.AddFileToZip(zipWriter, file)
		if err != nil {
			log.Println("Error adding file to zip:", err)
			context.JSON(http.StatusInternalServerError, "An internal error happened.")
			return
		}
	}
	zipWriter.Close()

	context.Header("Content-Disposition", "attachment; filename="+outputZip)
	context.Header("Content-Type", "application/zip")
	context.Header("Content-Length", strconv.Itoa(zipBuffer.Len()))
	context.Status(http.StatusOK)
	context.Writer.Write(zipBuffer.Bytes())
}

func checkGetTaskResult(context *gin.Context, nats_server *nats.Conn, queryParams *url.Values) bool {
	taskId := queryParams.Get(TASK_ID_PARAM)
	if taskId == "" {
		log.Println("Error: parameter taskId is missing")
		context.JSON(http.StatusBadRequest, gin.H{"error": "Parameter taskId is missing"})
		return false
	}

	userId := context.Request.Header.Get("X-Forwarded-User")

	log.Println("Received request to get result for task " + taskId)

	st, err := store.GetTaskStatus(nats_server, taskId, userId)
	switch err {
	case nil:
		break

	case errors.ErrUserNotFound, errors.ErrTaskNotFound:
		context.JSON(http.StatusBadRequest, gin.H{"error": "User does not have a task with given ID."})
		return false

	case errors.ErrUserInvalid:
		context.JSON(http.StatusUnauthorized, gin.H{"error": "User ID is invalid."})
		return false
	}

	if *st == task.PENDING || *st == task.EXECUTING {
		context.JSON(http.StatusOK, gin.H{"data": "Given task has not finished yet."})
		return false
	}

	return true
}

/************************ TASKS ************************/

func ReceiveTasks(nats_server *nats.Conn, handleFunc func(task.Task, *nats.Conn)) {
	err := SubscribeQueueTask(nats_server, handleFunc)
	if err != nil {
		log.Fatal("Unable to receive tasks. Aborting...")
	}
}

// Check PostTask request is valid
func checkPostTask(context *gin.Context, nats_server *nats.Conn, requestBody *request.PostTaskBody) bool {
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Request must have 'url' and 'parameters' parameters."})
		return false
	}

	if requestBody.Parameters == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'parameters' is missing."})
		return false
	}
	if requestBody.Url == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'url' is missing."})
		return false
	}

	user_id := context.Request.Header.Get("X-Forwarded-User")

	user_tasks, err := store.GetUserTasksId(nats_server, user_id)
	switch err {
	case nil, errors.ErrUserNotFound: // No errors or user has no tasks
		break

	case errors.ErrUserInvalid:
		context.JSON(http.StatusUnauthorized, gin.H{"error": "User ID is invalid."})
		return false

	default:
		context.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error has happened."})
		return false
	}

	max_requests_per_client, err := strconv.Atoi(os.Getenv("MAX_REQUESTS_PER_CLIENT"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error happened."})
		return false
	}
	if len(user_tasks) == max_requests_per_client {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "You already have the maximum number of requests allowed per client (" + os.Getenv("MAX_REQUESTS_PER_CLIENT") + ")."})
		return false
	}

	return true
}

func PostTask(context *gin.Context, nats_server *nats.Conn) {

	var requestBody request.PostTaskBody
	if !checkPostTask(context, nats_server, &requestBody) {
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

func GetAllTasks(context *gin.Context, nats_server *nats.Conn) {

	userId := context.Request.Header.Get("X-Forwarded-User")

	allTaskIds, err := store.GetUserTasksId(nats_server, userId)
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
