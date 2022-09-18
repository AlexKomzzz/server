package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func sayhello(c *gin.Context) {
	c.JSON(http.StatusOK, "Привет!")
}

func main() {

	mux := gin.New()

	mux.GET("/", sayhello)

	err := mux.Run(":8080") // устанавливаем порт веб-сервера
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
