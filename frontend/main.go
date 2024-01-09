// hemos instalado go get github.com/gin-gonic/gin
// he instalado go get github.com/gin-contrib/cors
package main

import (
	"ejemplo/clases"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var tareas = []clases.Tarea{}

func getTarea(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, tareas)
}

func addTarea(context *gin.Context) {
	var newTarea clases.Tarea

	if err := context.BindJSON(&newTarea); err != nil {
		return
	}

	tareas = append(tareas, newTarea)

	context.IndentedJSON(http.StatusCreated, newTarea)
}

func main() {
	//Creamos un servidor mediante el framework Gin
	router := gin.Default()

	//Configuramos el middleware CORS
	//para que puedan acceder al servidor desde un dominio diferente
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))

	//Para obtener los datos
	router.GET("/tareas", getTarea)
	//Para agregar datos
	router.POST("/tareas", addTarea)
	//Ejecutar el servidor
	router.Run(":80")
}
