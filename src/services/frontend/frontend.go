// hemos instalado go get github.com/gin-gonic/gin
// he instalado go get github.com/gin-contrib/cors
package main

import (
	. "cc/src/pkg/lib/TaskManager"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

func main() {
	//Creamos un servidor mediante el framework Gin
	router := gin.Default()

	//Configuramos el middleware CORS
	//para que puedan acceder al servidor desde un dominio diferente
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))

	nats_server, err := nats.Connect(os.Getenv("NATS_SERVER_ADDRESS")) //nats.DefaultURL
	if err != nil {
		log.Fatal(err)
	}
	defer nats_server.Close()

	//Para obtener los datos
	router.GET("/helloWorld", func(ctx *gin.Context) { ctx.IndentedJSON(http.StatusCreated, "Hola Mundo") })

	//Para obtener los datos
	router.GET("/getTaskResult", func(ctx *gin.Context) { GetTaskResult(ctx, nats_server) })
	//Para agregar datos
	router.POST("/createTask", func(ctx *gin.Context) { PostTask(ctx, nats_server) })

	router.GET("/getTaskStatus", func(ctx *gin.Context) { GetTaskStatus(ctx, nats_server) })

	router.GET("/getAllTasks", func(ctx *gin.Context) { GetAllTasks(ctx, nats_server) })

	router.GET("/metrics", func(c *gin.Context) {
		// Specify the new URL to redirect to
		newURL := "http://localhost:8080/metrics"
		// Perform the redirection
		c.Redirect(http.StatusFound, newURL)
	})
	//Ejecutar el servidor
	router.Run("0.0.0.0:8080")
}
