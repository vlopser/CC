package storemanager

import (
	"cc/src/pkg/models/errors"
	"cc/src/pkg/models/result"
	"cc/src/pkg/models/task"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

const (
	TTL_TASK = 5 * time.Minute
)

func SetTaskStatus(nats_server *nats.Conn, idUser string, idTask string, status task.Status) error {
	js, err := nats_server.JetStream()
	if err != nil {
		return err
	}

	user_bucket, err := js.KeyValue(idUser)
	switch err {
	case nil:
		break
	case nats.ErrBucketNotFound:
		user_bucket, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: idUser,
			TTL:    TTL_TASK,
		})
		switch err {
		case nil:
			break
		case nats.ErrInvalidBucketName:
			log.Println("User ID is invalid:", err.Error())
			return errors.ErrUserInvalid
		default:
			log.Println("Unexpected error:", err.Error())
			return err
		}
	default:
		log.Println("Unexpected error:", err.Error())
		return err
	}

	user_bucket.Put(idTask, []byte(fmt.Sprintf("%d", status)))
	switch err {
	case nil:
		break
	case nats.ErrInvalidKey:
		log.Println("Task id is invalid:", err.Error())
		return errors.ErrTaskInvalid
	default:
		log.Println("Unexpected error:", err.Error())
		return err
	}

	return nil
}

func GetTaskStatus(nats_server *nats.Conn, idTask string, idUser string) (*task.Status, error) {
	js, err := nats_server.JetStream()
	if err != nil {
		return nil, err
	}
	user_bucket, err := js.KeyValue(idUser)
	switch err {
	case nil:
		break

	case nats.ErrBucketNotFound:
		log.Println("User bucket does not exist:", err.Error())
		return nil, errors.ErrUserNotFound

	case nats.ErrInvalidBucketName:
		log.Println("User ID is invalid:", err.Error())
		return nil, errors.ErrUserInvalid

	default:
		log.Println("Unexpected error:", err.Error())
		return nil, err
	}

	kv_task_status, err := user_bucket.Get(idTask)
	switch err {
	case nil:
		break

	case nats.ErrKeyNotFound:
		log.Println("Task does not exist:", err.Error())
		return nil, errors.ErrTaskNotFound

	default:
		log.Println("Unexpected error:", err.Error())
		return nil, err
	}

	//Convert from []bytes to Status type
	status, _ := strconv.Atoi(string(kv_task_status.Value())) //Assume status was stored correctly as an Integer
	res := task.Status(status)

	return &res, nil

}

func GetUserTasksId(nats_server *nats.Conn, idUser string) ([]string, error) {
	js, err := nats_server.JetStream()
	if err != nil {
		return nil, err
	}
	user_bucket, err := js.KeyValue(idUser)
	switch err {
	case nil:
		break

	case nats.ErrBucketNotFound:
		log.Println("User bucket does not exist:", err.Error())
		return nil, errors.ErrUserNotFound

	case nats.ErrInvalidBucketName:
		log.Println("User ID is invalid:", err.Error())
		return nil, errors.ErrUserInvalid

	default:
		log.Println("Unexpected error:", err.Error())
		return nil, err
	}

	ch_user_tasks, err := user_bucket.ListKeys()
	if err != nil {
		log.Println("Unexpected error:", err.Error())
		return nil, err
	}

	result := make([]string, 0)

	for task := range ch_user_tasks.Keys() {
		result = append(result, task)
	}

	return result, nil

}

func StoreFileInBucket(nats_server *nats.Conn, file_path string, file_name string, bucket_name string) error {
	js, err := nats_server.JetStream()
	if err != nil {
		return err
	}

	bucket, err := js.ObjectStore(bucket_name)
	if err != nil {
		return err
	}

	// We override behaviout of PutFile so we can choose arguments as filename
	f, err := os.Open(file_path)
	if err != nil {
		return err
	}
	defer f.Close()

	bucket.Put(&nats.ObjectMeta{Name: file_name}, f)
	if err != nil {
		return err
	}

	return nil
}

func CreateTaskBucket(nats_server *nats.Conn, task_id string) error {
	js, err := nats_server.JetStream()
	if err != nil {
		return err
	}

	_, err = js.CreateObjectStore(&nats.ObjectStoreConfig{
		Bucket:   task_id,
		TTL:      TTL_TASK, //Time until bucket is automatically deleted
		MaxBytes: 10000,    //Only keep 10 MB maximum
	})
	switch err {
	case nil:
		break
	case nats.ErrInvalidStoreName:
		log.Println("Task ID is invalid:", err.Error())
		return errors.ErrTaskInvalid
	default:
		log.Println("Unexpected error:", err.Error())
		return err
	}

	return nil
}

func GetResult(nats_server *nats.Conn, taskId string) (*result.Result, error) {
	js, err := nats_server.JetStream()
	if err != nil {
		return nil, err
	}

	task_bucket, err := js.ObjectStore(taskId)
	switch err {
	case nil:
		break
	case nats.ErrInvalidStoreName:
		log.Println("Task ID is invalid:", err.Error())
		return nil, errors.ErrTaskInvalid

	case nats.ErrStreamNotFound:
		log.Println("Task bucket not found:", err.Error())
		return nil, errors.ErrTaskNotFound

	default:
		log.Println("Unexpected error:", err.Error())
		return nil, err
	}

	parsedTaskId, err := uuid.Parse(taskId)
	if err != nil {
		fmt.Println("Error parsing UUID:", err)
		return nil, err
	}

	res := result.Result{
		TaskId: parsedTaskId,
	}

	files_in_bucket, _ := task_bucket.List()
	for _, file := range files_in_bucket {
		res_filename := res.TaskId.String() + "_" + file.Name
		task_bucket.GetFile(file.Name, res_filename)
		res.Files = append(res.Files, res_filename)
	}

	return &res, nil
}
