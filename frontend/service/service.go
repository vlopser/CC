package service

import (
	"frontend/classes"
	"frontend/natsUtils"
	"github.com/gin-gonic/gin"
	"gopkg.in/src-d/go-git.v4"
	"log"
	"math/rand"
	"net/http"
)

func checkGitRepo(url string) {

	// Attempt to open the repository
	_, err := git.PlainOpen(url)
	if err == nil {
		log.Println("The directory is a Git repository.")
	} else if err == git.ErrRepositoryNotExists {
		log.Println("The directory is not a Git repository.")
	}
}

func GetResult(context *gin.Context) {
	// tenemos que tomar los resultados desde un object store
}

func PostTask(context *gin.Context) {

	// open connection to nats
	conn := natsUtils.GetConnection()
	defer conn.Close()

	var requestBody classes.RequestBody
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	log.Println("Received request to create a task from " + requestBody.Url)
	log.Println("Task parameters are:")
	for _, param := range requestBody.Parameters {
		log.Println(param)
	}

	checkGitRepo(requestBody.Url)

	task := classes.NewTask(rand.Int(), requestBody.Url, requestBody.Parameters)
	//natsUtils.Publish(conn, &task)

	context.IndentedJSON(http.StatusCreated, task.IdTask)
}
