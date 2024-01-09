package storemanager

import (
	models "cc/Models"
	"fmt"

	"github.com/nats-io/nats.go"
)

const (
	STATUS_BUCKET  = "status_bucket"
	RESULTS_BUCKET = "results_bucket"
)

func ChangeState(nats_server *nats.Conn, idTask string, status models.Status) error {
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

func StoreResult(nats_server *nats.Conn, result models.Result) error {
	js, err := nats_server.JetStream()
	if err != nil {
		return err
	}

	results_bucket, err := js.ObjectStore(RESULTS_BUCKET)
	if err != nil {
		return err
	}

	_, err = results_bucket.PutString(result.TaskId.String(), result.Output)
	if err != nil {
		return err
	}

	return nil
}

func GetResult(nats_server *nats.Conn, taskId string) (*models.Result, error) {
	js, err := nats_server.JetStream()
	if err != nil {
		return nil, err
	}

	results_bucket, err := js.ObjectStore(RESULTS_BUCKET)
	if err != nil {
		return nil, err
	}

	var result models.Result
	result.Output, err = results_bucket.GetString(taskId)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
