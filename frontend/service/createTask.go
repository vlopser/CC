package service

import (
	"frontend/classes"
	"github.com/gin-gonic/gin"
	"gopkg.in/src-d/go-git.v4"
	"log"
	"math/rand"
	"net/http"
)

func PostTask(context *gin.Context) {

	var requestBody classes.PostTaskBody
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Received request to create a task from " + requestBody.Url)
	log.Println("Task parameters are:")
	for _, param := range requestBody.Parameters {
		log.Println(param)
	}

	_, err := git.PlainOpen(requestBody.Url)
	if err == nil {
		log.Println("Error: The directory is a Git repository.", err)
	} else if err == git.ErrRepositoryNotExists {
		// definimos un codigo de error para errores genericos
		log.Println("Error: The directory is not a Git repository.", err)
		context.IndentedJSON(1100, nil)
		return
	}

	task := classes.Task{
		IdTask:     rand.Int(),
		RepoUrl:    requestBody.Url,
		Parameters: requestBody.Parameters,
		Status:     100,
	}

	// todo llamar la libreria
	context.IndentedJSON(http.StatusCreated, task.IdTask)
}
