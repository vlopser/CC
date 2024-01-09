// hemos instalado go get github.com/gin-gonic/gin
// he instalado go get github.com/gin-contrib/cors
package main

import (
	"frontend/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	//Creamos un servidor mediante el framework Gin
	router := gin.Default()

	//Configuramos el middleware CORS
	//para que puedan acceder al servidor desde un dominio diferente
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))

	//Para obtener los datos
	router.GET("/getResult", service.GetResult)
	//Para agregar datos
	router.POST("/createTask", service.PostTask)
	//Ejecutar el servidor
	router.Run(":8080")
}
