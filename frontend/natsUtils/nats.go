package natsUtils

import (
	"encoding/json"
	"frontend/classes"
	"github.com/nats-io/nats.go"
	"log"
)

func GetConnection() *nats.Conn {
	// Connect to the NATS server
	conn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func Publish(conn *nats.Conn, task *classes.Task) {

	pkgJSON, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Publish("putTask", pkgJSON)
	if err != nil {
		log.Fatal(err)
	}
}
