package service

import (
	"frontend/classes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetResult(context *gin.Context) {

	var requestBody *classes.GetResultBody
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Received request to get result for task " + requestBody.TaskId)

	//conn, err := nats.Connect("nats://localhost:4222")
	//if err != nil {
	//	log.Println("It was impossible to open connection to nats queue", err)
	//	context.IndentedJSON(1100, nil)
	//	return
	//}
	//
	//defer conn.Close()
	//
	//js, err := conn.JetStream()
	//if err != nil {
	//	context.IndentedJSON(1100, nil)
	//}
	//
	//// todo define bucket name
	//resultBucket, err := js.ObjectStore("")
	//if err != nil {
	//	context.IndentedJSON(1100, nil)
	//	return
	//}
	//
	//// todo define bucket name
	//var result classes.Result
	//result.Output, err = resultBucket.GetString(requestBody.TaskId)
	//if err != nil {
	//	context.IndentedJSON(1100, nil)
	//	return
	//}

	//if result.Output == "" {
	//	// si el output esta vacio la execucion del task no ha terminado todavia o ha pasado algo?
	//	log.Println("The result for the task is not available")
	//	context.IndentedJSON(http.StatusOK, "The result is not available")
	//	return
	//}

	// todo llamar a la libreria

	context.IndentedJSON(http.StatusCreated, "ok")
}
