package service

import (
	"frontend/classes"
	"frontend/natsUtils"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"gopkg.in/src-d/go-git.v4"
	"log"
	"math/rand"
	"net/http"
)

func PostTask(context *gin.Context) {

	// open connection to nats
	conn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Println("It was impossible to open connection to nats queue", err)
		context.IndentedJSON(1100, nil)
		return
	}

	defer conn.Close()

	var requestBody classes.PostTaskBody
	if err = context.ShouldBindJSON(&requestBody); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Received request to create a task from " + requestBody.Url)
	log.Println("Task parameters are:")
	for _, param := range requestBody.Parameters {
		log.Println(param)
	}

	_, err = git.PlainOpen(requestBody.Url)
	if err == nil {
		log.Println("The directory is a Git repository.", err)
	} else if err == git.ErrRepositoryNotExists {
		// definimos un codigo de error para errores genericos
		context.IndentedJSON(1100, nil)
		return
	}

	task := classes.Task{
		IdTask:     rand.Int(),
		RepoUrl:    requestBody.Url,
		Parameters: requestBody.Parameters,
		Status:     100,
	}

	// todo define subject name
	err = natsUtils.Publish("", conn, &task)
	if err != nil {
		context.IndentedJSON(1100, nil)
		return
	}

	context.IndentedJSON(http.StatusCreated, task.IdTask)
}
