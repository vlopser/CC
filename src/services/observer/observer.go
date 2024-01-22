package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
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

const (
	REQUEST_QUEUE = "request_queue"
)

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
	ticker1 := time.NewTicker(time.Second * 2)
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

func getMetrics(ctx *gin.Context) {
	//endpoint := "http://localhost:8222/accstatz"
	// todo put other metrics that you consider important

}

func main() {
	router := gin.Default()

	natsServer, err := nats.Connect(os.Getenv("NATS_SERVER_ADDRESS")) //nats.DefaultURL
	if err != nil {
		log.Fatal(err)
	}
	defer natsServer.Close()

	router.POST("/metrics", func(ctx *gin.Context) { getMetrics(ctx) })
	router.Run(":8080")
	go scheduledThread()
}
