package natsUtils

import (
	"encoding/json"
	"frontend/classes"
	"github.com/nats-io/nats.go"
	"log"
	"testing"
	"time"
)

type NATSQueue struct {
	conn *nats.Conn
	sub  *nats.Subscription
}

func NewNATSQueue(subject string) (*NATSQueue, error) {
	conn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("It was impossible to open connection to nats queue", err)
	}
	sub, err := conn.QueueSubscribe(subject, "test", func(msg *nats.Msg) {
		log.Println("Received message")
		var task classes.Task
		json.Unmarshal(msg.Data, &task)
		log.Printf("Received task with id %d", task.IdTask)
	})
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &NATSQueue{conn: conn, sub: sub}, nil
}

func TestPublish(t *testing.T) {
	queue, err := NewNATSQueue("putTaskTest")
	if err != nil {
		log.Fatal("Impossible to create mocked nats queue")
	}
	type args struct {
		conn *nats.Conn
		task *classes.Task
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "testPutTask",
			args: args{
				conn: queue.conn,
				task: &classes.Task{
					IdTask:     1,
					RepoUrl:    "test",
					Parameters: nil,
					Status:     100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Println("Start testPutTask test...")
			Publish("putTaskTest", tt.args.conn, tt.args.task)
			time.Sleep(time.Second * 2)
			queue.conn.Close()
		})
	}
}
