package storemanager

import (
	"cc/pkg/models/result"
	"cc/pkg/models/task"
	"fmt"

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

func StoreResult(nats_server *nats.Conn, result result.Result) error {
	js, err := nats_server.JetStream()
	if err != nil {
		return err
	}

	results_bucket, err := js.ObjectStore(RESULTS_BUCKET)
	if err != nil {
		println("ey")
		return err
	}

	_, err = results_bucket.PutString(result.TaskId.String(), result.Output)
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

	results_bucket, err := js.ObjectStore(RESULTS_BUCKET)
	if err != nil {
		return nil, err
	}

	var result result.Result
	result.Output, err = results_bucket.GetString(taskId)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
