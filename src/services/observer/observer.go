package main

import (
	storemanager "cc/src/pkg/lib/StoreManager"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"strconv"
	"time"
)

var count int
var limit = 10
var workers = 1
var natsServer *nats.Conn
var congested = false

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
	inMsgs, outMsgs, err := storemanager.GetInOutMsgs(natsServer)
	if err != nil {
		log.Println("An error occurred getting in/out msgs from observer KV", err.Error())
		return
	} else if "" == inMsgs || "" == outMsgs {
		return
	}

	in, err := strconv.Atoi(inMsgs)
	if err != nil {
		log.Println("Unexpected error:", err.Error())
		return
	}
	out, err := strconv.Atoi(outMsgs)
	if err != nil {
		log.Println("Unexpected error:", err.Error())
		return
	}

	log.Printf("Input messages: %d", in)
	log.Printf("Output messages: %d", out)

	difference := in - out
	key, event := decideNumberOfWorkers(difference)
	err = storemanager.SetSystemStatus(natsServer, congested, workers)
	if err != nil {
		log.Printf("Error creating event: %s %s", event, err.Error())
	}
	if key != "" && event != "" {
		err = storemanager.CreateObserverEvent(natsServer, key, event)
		if err != nil {
			log.Printf("Error creating event: %s %s", event, err.Error())
		}
	}
}

func decideNumberOfWorkers(difference int) (string, string) {
	if workers > 1 && difference < limit/3 {
		currentTime := time.Now()
		timeString := currentTime.Format("2006-01-02 15:04:05")
		workers--
		aLog := fmt.Sprintf("%s Difference between sent and received message is %d. Reducing number of workers, n° workers : %d", timeString, difference, workers)
		log.Printf(aLog)
		congested = false
		return currentTime.Format("150405"), aLog
	} else if difference > limit {
		if count == 3 {
			currentTime := time.Now()
			timeString := currentTime.Format("2006-01-02 15:04:05")
			workers++
			aLog := fmt.Sprintf("%s Nats queue has been busy for the last three times. The system is congested. Difference between sent and received message is %d. Adding new worker, n° workers : %d", timeString, difference, workers)
			log.Printf(aLog)
			count = 0
			congested = true
			return currentTime.Format("150405"), aLog
		} else {
			currentTime := time.Now()
			timeString := currentTime.Format("2006-01-02 15:04:05")
			aLog := fmt.Sprintf("%s Nats queue is busy, system is getting in congestion. Difference between sent and received message is %d", timeString, difference)
			log.Println(aLog)
			count++
			return currentTime.Format("150405"), aLog
		}
	} else {
		congested = false
		return "", ""
	}
}

func main() {
	var err error
	natsServer, err = nats.Connect(os.Getenv("NATS_SERVER_ADDRESS")) //nats.DefaultURL
	if err != nil {
		log.Fatal(err)
	}
	defer natsServer.Close()

	scheduledThread()

}
