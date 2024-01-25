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
	TTL_TASK      = 5 * time.Minute
	TTL_LOGS      = 24 * time.Hour
	OBSERVER_KV   = "observer"
	SYSTEM_STATUS = "system_status"
	WORKERS       = "workers"
	IN_MSGS       = "in_msgs"
	OUT_MSGS      = "out_msgs"
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
		user_bucket, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: idUser,
			TTL:    TTL_TASK,
		})
		switch err {
		case nil:
			break
		case nats.ErrInvalidBucketName:
			log.Println("User ID is invalid:", err.Error())
			return nil, errors.ErrUserInvalid
		default:
			log.Println("Unexpected error:", err.Error())
			return nil, err
		}

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

func CreateObserverEvent(nats_server *nats.Conn, key string, event string) error {
	js, err := nats_server.JetStream()
	if err != nil {
		return err
	}

	bucket, err := js.KeyValue(OBSERVER_KV)
	switch err {
	case nil:
		break
	case nats.ErrBucketNotFound:
		log.Println("Observer bucket does not exist:", err.Error())
		return errors.ErrUserNotFound
	default:
		log.Println("Unexpected error:", err.Error())
		return err
	}
	_, err = bucket.Create(key, []byte(event))
	switch err {
	case nil:
		break
	case nats.ErrInvalidKey:
		log.Println("Bucket key is invalid:", err.Error())
		return errors.ErrTaskInvalid
	default:
		log.Println("Unexpected error:", err.Error())
		return err
	}

	return nil
}

func GetObserverLogs(nats_server *nats.Conn) ([]string, error) {
	js, err := nats_server.JetStream()
	if err != nil {
		return nil, err
	}
	bucket, err := js.KeyValue(OBSERVER_KV)
	switch err {
	case nil:
		break

	case nats.ErrBucketNotFound:
		log.Println("Observer bucket does not exist:", err.Error())
		return nil, errors.ErrUserNotFound

	case nats.ErrInvalidBucketName:
		log.Println("Bucket name is invalid:", err.Error())
		return nil, errors.ErrUserInvalid

	default:
		log.Println("Unexpected error:", err.Error())
		return nil, err
	}

	logs, err := bucket.ListKeys()
	if err != nil {
		log.Println("Unexpected error:", err.Error())
		return nil, err
	}

	result := make([]string, 0)

	for event := range logs.Keys() {
		if event == WORKERS || event == SYSTEM_STATUS || event == IN_MSGS || event == OUT_MSGS {
			continue
		}
		value, err := bucket.Get(event)
		if err != nil {
			log.Println(err.Error())
		}
		result = append(result, string(value.Value()))
	}

	return result, nil
}

func SetSystemStatus(natsServer *nats.Conn, congestionated bool, nWorkers int) error {
	js, err := natsServer.JetStream()
	if err != nil {
		return err
	}

	bucket, err := js.KeyValue(OBSERVER_KV)
	switch err {
	case nil:
		break
	case nats.ErrBucketNotFound:
		log.Println("Observer bucket does not exist:", err.Error())
		return errors.ErrUserNotFound
	default:
		log.Println("Unexpected error:", err.Error())
		return err
	}

	if congestionated {
		_, err = bucket.Put(SYSTEM_STATUS, []byte("congestionated"))
	} else {
		_, err = bucket.Put(SYSTEM_STATUS, []byte("OK"))
	}
	switch err {
	case nil:
		break
	case nats.ErrInvalidKey:
		log.Println("Bucket key is invalid:", err.Error())
		return errors.ErrTaskInvalid
	default:
		log.Println("Unexpected error:", err.Error())
		return err
	}

	_, err = bucket.Put(WORKERS, []byte(strconv.Itoa(nWorkers)))
	switch err {
	case nil:
		break
	case nats.ErrInvalidKey:
		log.Println("Bucket key is invalid:", err.Error())
		return errors.ErrTaskInvalid
	default:
		log.Println("Unexpected error:", err.Error())
		return err
	}
	return nil
}

func GetSystemStatus(nats_server *nats.Conn) (string, string, error) {

	js, err := nats_server.JetStream()
	if err != nil {
		return "", "", err
	}
	bucket, err := js.KeyValue(OBSERVER_KV)
	switch err {
	case nil:
		break

	case nats.ErrBucketNotFound:
		log.Println("Observer bucket does not exist:", err.Error())
		return "", "", errors.ErrUserNotFound

	case nats.ErrInvalidBucketName:
		log.Println("Bucket name is invalid:", err.Error())
		return "", "", errors.ErrUserInvalid

	default:
		log.Println("Unexpected error:", err.Error())
		return "", "", err
	}

	status, err := bucket.Get(SYSTEM_STATUS)
	if err != nil {
		log.Println(err.Error())
	}
	nWorkers, err := bucket.Get(WORKERS)
	if err != nil {
		log.Println(err.Error())
	}

	return string(status.Value()), string(nWorkers.Value()), nil
}

func SetInOutMsgs(natsServer *nats.Conn, key string) error {
	js, err := natsServer.JetStream()
	if err != nil {
		return err
	}

	bucket, err := js.KeyValue(OBSERVER_KV)
	switch err {
	case nil:
		break
	case nats.ErrBucketNotFound:
		log.Println("Observer bucket does not exist:", err.Error())
		return errors.ErrUserNotFound
	default:
		log.Println("Unexpected error:", err.Error())
		return err
	}

	// get last number of sent/received message
	value, err := bucket.Get(key)
	switch err {
	case nil:
		break
	case nats.ErrKeyNotFound:
		_, err = bucket.Create(key, []byte(strconv.Itoa(0)))
		if err != nil {
			log.Println("Unexpected error:", err.Error())
			return err
		}
		break
	default:
		log.Println("Error retrieving number of "+key, err.Error())
		return err
	}

	inMsgs, err := strconv.Atoi(string(value.Value()))
	if err != nil {
		log.Println("Unexpected error:", err.Error())
		return err
	}
	_, err = bucket.Put(key, []byte(strconv.Itoa(inMsgs+1)))

	switch err {
	case nil:
		break
	case nats.ErrInvalidKey:
		log.Println("Bucket key is invalid:", err.Error())
		return errors.ErrTaskInvalid
	default:
		log.Println("Unexpected error:", err.Error())
		return err
	}

	return nil
}

func GetInOutMsgs(natsServer *nats.Conn) (string, string, error) {
	js, err := natsServer.JetStream()
	if err != nil {
		return "", "", err
	}
	bucket, err := js.KeyValue(OBSERVER_KV)
	switch err {
	case nil:
		break

	case nats.ErrBucketNotFound:
		bucket, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: OBSERVER_KV,
			TTL:    TTL_LOGS,
		})
		switch err {
		case nil:
			break
		case nats.ErrInvalidBucketName:
			log.Println("Bucket name is invalid:", err.Error())
			return "", "", errors.ErrUserInvalid
		default:
			log.Println("Unexpected error:", err.Error())
			return "", "", err
		}
	case nats.ErrInvalidBucketName:
		log.Println("Bucket name is invalid:", err.Error())
		return "", "", errors.ErrUserInvalid

	default:
		log.Println("Unexpected error:", err.Error())
		return "", "", err
	}

	inMsgs, err := bucket.Get(IN_MSGS)
	switch err {
	case nil:
		break
	case nats.ErrKeyNotFound:
		_, err = bucket.Create(IN_MSGS, []byte(strconv.Itoa(0)))
		return "", "", err
	default:
		log.Println("Error retrieving number of "+IN_MSGS, err.Error())
		return "", "", err
	}

	outMsgs, err := bucket.Get(OUT_MSGS)
	switch err {
	case nil:
		break
	case nats.ErrKeyNotFound:
		_, err = bucket.Create(OUT_MSGS, []byte(strconv.Itoa(0)))
		return "", "", err
	default:
		log.Println("Error retrieving number of "+OUT_MSGS, err.Error())
		return "", "", err
	}

	return string(inMsgs.Value()), string(outMsgs.Value()), nil
}
