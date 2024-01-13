// hemos instalado go get github.com/gin-gonic/gin
// he instalado go get github.com/gin-contrib/cors
package main

import (
	. "cc/pkg/lib/TaskManager"

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

	// nats_server, err := nats.Connect(nats.DefaultURL)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer nats_server.Close()

	var nats_server *nats.Conn

	//Para obtener los datos
	router.GET("/getResult", func(ctx *gin.Context) { GetTaskResult(ctx, nats_server) })
	//Para agregar datos
	router.POST("/createTask", func(ctx *gin.Context) { CreateTask(ctx, nats_server) })
	//Ejecutar el servidor
	router.Run("localhost:8080")
}
