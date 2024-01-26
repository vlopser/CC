package main

import (
	injectionmanager "cc/src/pkg/lib/InjectionManager"
	storemanager "cc/src/pkg/lib/StoreManager"
	"cc/src/pkg/models/requestInjection"
	"cc/src/pkg/utils"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	INJECTOR_BUCKET = "injector_bucket"
)

func waitForInjectionRequests(nats_server *nats.Conn, wg *sync.WaitGroup) {
	injectionmanager.ReceiveRequest(nats_server, handleRequest)
	wg.Add(1)

	log.Println("Waiting for injection requests. (Press Ctrl+C to exit)")

	wg.Wait()
	time.Sleep(1 * time.Second)
}

func handleRequest(nats_server *nats.Conn, req requestInjection.RequestInjection) {
	loop := true
	for loop {
		err := storemanager.StoreStringInBucket(nats_server, req.File_content, req.File_name, INJECTOR_BUCKET)
		switch err {
		case nil:
			log.Println("File successfully injected")
			loop = false
			break

		case nats.ErrStreamNotFound:
			log.Println("Creating injector bucket")
			storemanager.CreateInjectorBucket(nats_server, INJECTOR_BUCKET)

		default:
			log.Fatal("Unexpected error happended:", err.Error())
		}
	}
}

func main() {

	nats_server, err := nats.Connect(os.Getenv("NATS_SERVER_ADDRESS"))
	if err != nil {
		log.Fatal(err)
	}
	defer nats_server.Close()

	wg := sync.WaitGroup{}
	go utils.WaitForSigkill(&wg)

	waitForInjectionRequests(nats_server, &wg)

}
