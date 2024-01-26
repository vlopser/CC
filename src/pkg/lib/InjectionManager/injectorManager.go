package injectionmanager

import (
	. "cc/src/pkg/lib/QueueManager"
	queuemanager "cc/src/pkg/lib/QueueManager"
	"cc/src/pkg/models/requestInjection"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

func ReceiveRequest(nats_server *nats.Conn, handleFunc func(*nats.Conn, requestInjection.RequestInjection)) {
	err := SubscribeQueueInjection(nats_server, handleFunc)
	if err != nil {
		log.Fatal("Unable to receive injection requests. Aborting...")
	}
}

func PostInjection(context *gin.Context, nats_server *nats.Conn) {

	requester_mail := context.Request.Header.Get("X-Forwarded-Email")
	if !strings.Contains(os.Getenv("ADMIN_USERS"), requester_mail) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Only admin users are allowed to inject files."})
		return
	}

	file, header, err := context.Request.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No file given."})
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error when reading file"})
		return
	}

	queuemanager.EnqueueInjectionRequest(
		requestInjection.RequestInjection{
			File_name:    header.Filename,
			File_content: string(content),
		}, nats_server)

	// for key, values := range context.Request.Header {
	// 	for _, value := range values {
	// 		log.Printf("%s: %s\n", key, value)
	// 	}
	// }
}
