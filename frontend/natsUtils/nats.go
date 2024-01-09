package natsUtils

import (
	"encoding/json"
	"frontend/classes"
	"github.com/nats-io/nats.go"
	"log"
)

func Publish(subject string, conn *nats.Conn, task *classes.Task) error {

	pkgJSON, err := json.Marshal(task)
	if err != nil {
		log.Println("An error occurred while trying to marshalling json input", err)
	}
	err = conn.Publish(subject, pkgJSON)
	if err != nil {
		log.Println("An error occurred in publishing task to nats queue", err)
	}
	return err
}
