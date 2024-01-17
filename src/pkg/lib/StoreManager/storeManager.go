package storemanager

import (
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
	STATUS_BUCKET  = "status_bucket"
	RESULTS_BUCKET = "results_bucket"
)

func ChangeState(nats_server *nats.Conn, idTask string, idUser string, status task.Status) error {
	js, err := nats_server.JetStream()
	if err != nil {
		return err
	}

	user_bucket, err := js.KeyValue(idUser)
	if err != nil {
		log.Println(err.Error())
		user_bucket, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: idUser,
		})
		if err != nil {
			log.Println(err.Error())
			return err
		}
		// return err
	}

	user_bucket.Put(idTask, []byte(fmt.Sprintf("%d", status)))
	if err != nil {
		return err
	}

	return nil
}

func GetTaskState(nats_server *nats.Conn, idTask string, idUser string) (*task.Status, error) {
	js, err := nats_server.JetStream()
	if err != nil {
		return nil, err
	}
	user_bucket, err := js.KeyValue(idUser)
	if err != nil {
		log.Println("User bucket does not exist:", err.Error())
		return nil, err
	}

	task_status, err := user_bucket.Get(idTask)
	if err != nil {
		log.Println("Task KV does not exist:", err.Error())
		return nil, err
	}

	//Convert from []bytes to Status type
	val, _ := strconv.Atoi(string(task_status.Value()))
	res := task.Status(val)

	return &res, nil

}

func GetAllTaskId(nats_server *nats.Conn, idUser string) ([]string, error) {
	js, err := nats_server.JetStream()
	if err != nil {
		return nil, err
	}
	user_bucket, err := js.KeyValue(idUser)
	if err != nil {
		log.Println("User bucket does not exist:", err.Error())
		return nil, err
	}

	allTasks, _ := user_bucket.ListKeys()
	// if err != nil {
	// 	log.Println("")
	// }

	// allTaskIds := <-allTasks.Keys()

	result := make([]string, 0)

	for key := range allTasks.Keys() {
		result = append(result, key)
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
		TTL:      5 * time.Minute, //Time until bucket is automatically deleted
		MaxBytes: 10000,           //Only keep 10 MB maximum
	})
	if err != nil {
		return err
	}

	return nil
}

func GetResult(nats_server *nats.Conn, taskId string) (*result.Result, error) {
	js, err := nats_server.JetStream()
	if err != nil {
		return nil, err
	}

	bucket, err := js.ObjectStore(taskId)
	if err != nil {
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

	files_in_bucket, _ := bucket.List()
	for _, file := range files_in_bucket {
		bucket.GetFile(file.Name, "frontend_"+file.Name)
		res.Files = append(res.Files, file.Name)
	}

	return &res, nil
}
