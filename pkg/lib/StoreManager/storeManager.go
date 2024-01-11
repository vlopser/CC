package storemanager

import (
	"cc/pkg/models/result"
	"cc/pkg/models/task"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

const (
	STATUS_BUCKET  = "status_bucket"
	RESULTS_BUCKET = "results_bucket"
)

func ChangeState(nats_server *nats.Conn, idTask string, status task.Status) error {
	js, err := nats_server.JetStream()
	if err != nil {
		return err
	}

	status_bucket, err := js.KeyValue(STATUS_BUCKET)
	if err != nil {
		return err
	}

	status_bucket.Put(idTask, []byte(fmt.Sprintf("%d", status)))
	if err != nil {
		return err
	}

	return nil
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

	file_output := "output_" + taskId
	bucket.GetFile(task.OUTPUT_FILE, file_output)
	if err != nil {
		return nil, err
	}

	file_errors := "errors_" + taskId
	bucket.GetFile(task.ERRORS_FILE, file_errors)
	if err != nil {
		return nil, err
	}

	//TODO: Check if error file is empty, dont return it

	parsedUUID, err := uuid.Parse(taskId)
	if err != nil {
		fmt.Println("Error parsing UUID:", err)
		return nil, err
	}
	res := result.Result{
		TaskId: parsedUUID,
		Output: file_output,
		Errors: file_errors,
	}

	return &res, nil
}
