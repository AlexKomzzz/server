package handler

import "github.com/gin-gonic/gin"

func InitRouter() {
	mux := gin.New()
	mux.Group("/auth")
	{
		mux.POST("/sign-up", signUp)
		mux.POST("/sign-in", signIn)
	}
}
