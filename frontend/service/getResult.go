package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetResult(context *gin.Context) {
	queryParams := context.Request.URL.Query()
	taskId := queryParams.Get("taskId")
	if taskId == "" {
		log.Println("Error: parameter taskId is missing")
		context.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	log.Println("Received request to get result for task " + taskId)

	// todo llamar a la libreria

	context.IndentedJSON(http.StatusCreated, "ok")
}
