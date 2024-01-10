package service

import (
	"frontend/classes"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/src-d/go-git.v4"
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

	tempFolder, err := os.MkdirTemp("", "tmp")
	if err != nil {
		log.Println("Error creating temp folder")
		context.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}

	// check git repo
	log.Println("Trying to clone git repo to check is validity...")
	_, err = git.PlainClone(tempFolder, false, &git.CloneOptions{
		URL: requestBody.Url,
	})
	if err == nil {
		log.Println("The url represent a Git repository.")
	} else if err != nil {
		log.Println("Error: The directory is not a valid Git repository.", err)
		context.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	err = os.RemoveAll(tempFolder)
	if err != nil {
		log.Println("Error trying to delete temp folder")
		context.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}

	task := classes.Task{
		TaskId:     rand.Int(),
		UserMail:   context.Request.Header.Get("X-Forwarded-Email"),
		RepoUrl:    requestBody.Url,
		Parameters: requestBody.Parameters,
		Status:     100,
	}

	// todo llamar la libreria
	context.IndentedJSON(http.StatusCreated, task)
}
