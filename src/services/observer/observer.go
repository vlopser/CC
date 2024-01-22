package main

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var count int
var lowerLimit = 200
var upperLimit = 1000
var workers = 1

type body struct {
	ServerId     string         `json:"server_id"`
	Now          time.Time      `json:"now"`
	AccountStatz []accountStatz `json:"account_statz"`
}

type sentReceived struct {
	Msgs  int `json:"msgs"`
	Bytes int `json:"bytes"`
}

type accountStatz struct {
	Acc              string       `json:"acc"`
	Conns            int          `json:"conns"`
	Leafnodes        int          `json:"leafnodes"`
	TotalConns       int          `json:"total_conns"`
	NumSubscriptions int          `json:"num_subscriptions"`
	Sent             sentReceived `json:"sent"`
	Received         sentReceived `json:"received"`
	SlowConsumers    int          `json:"slow_consumers"`
}

var running = false

func scheduledThread() {
	ticker1 := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker1.C:
			checkQueue()
		}
	}
}

func checkQueue() {
	endpoint := "http://localhost:8222/accstatz"
	response, err := http.Post(endpoint, "application/json", nil)
	if err != nil {
		log.Printf("Error connecting to endpoint %s", endpoint)
	}
	var msg body
	bodyText, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading JSON:", err)
		return
	}
	err = json.Unmarshal(bodyText, &msg)
	if err != nil {
		log.Println("Error when unmarshalling JSON:", err)
		return
	}
	log.Println("Response is:\n", msg)

	totalMsgs := msg.AccountStatz[0].Sent.Msgs + msg.AccountStatz[0].Received.Msgs - count
	if workers > 1 && totalMsgs < lowerLimit {
		log.Printf("Received/sent messages : %d", totalMsgs)
		log.Printf("Limit : %d", lowerLimit)
		workers--
		log.Printf("Reducing number of workers, n° workers : %d", workers)
		upperLimit /= 2 * workers
	} else if totalMsgs > upperLimit {
		log.Println("Nats queue has exceeded current limit of sent/received messages")
		log.Printf("Received/sent messages : %d", totalMsgs)
		log.Printf("Limit : %d", upperLimit)
		workers++
		log.Printf("Adding new worker, n° workers : %d", workers)
		upperLimit *= 2
	}
}

func main() {
	/*
		// Replace "nats://your-nats-server-url:4222" with the actual URL of your NATS server
		natsURL := "nats://your-nats-server-url:4222"
		queueName := "your_queue_name"

		nc, err := nats.Connect(natsURL)
		if err != nil {
			log.Fatalf("Error connecting to NATS server: %v", err)
		}
		defer nc.Close()

		queueInfo, err := nc.Request(fmt.Sprintf("nats.queue.%s", queueName), nil, time.Second)
		if err != nil {
			log.Fatalf("Error requesting queue info: %v", err)
		}

		fmt.Printf("Queue Info: %s\n", queueInfo.Data)
	*/
	nats_server, err := nats.Connect(os.Getenv("NATS_SERVER_ADDRESS")) //nats.DefaultURL
	if err != nil {
		log.Fatal(err)
	}
	defer nats_server.Close()

	go scheduledThread()
}
