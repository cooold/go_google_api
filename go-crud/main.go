package main

import (
	"github.com/Qbom/go-crud/config"
	"github.com/Qbom/go-crud/controllers"
	"github.com/Qbom/go-crud/handler"
	"github.com/Qbom/go-crud/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	config.Init()
}

func RequestIDMiddleware(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	//ctx.Header("Content-Type", "application/json;charset=utf-8")
	//ctx.Header("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI4NTNlODUwNmFiZWQ5MmM5ZTNmMWM4Yzk1ZTdhZmU4ZCIsInN1YiI6IjYzOGI0NDcwN2Q1ZGI1MGZiZjAwNjFjZCIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.aC3CaN69mIPG-NCJChRUV7LPxkXRUrxLxMk6Yfv64NA")
	ctx.Next()
}

func main() {
	r := gin.Default()
	r.Use(RequestIDMiddleware)
	r.GET("/TMDBList/:page", controllers.TMDBList)
	r.POST("/posts", controllers.PostsCreate)
	r.PUT("/posts/:id", controllers.PostUpdate)
	r.GET("/posts", controllers.PostsIndex)
	r.GET("/posts/:id", controllers.PostsShow)
	r.DELETE("/posts/:id", controllers.PostDelet)

	api := r.Group("/api")
	{
		api.GET("ouath/google/url", handler.GoogleAccsess)
		api.GET("ouath/google/login", handler.GoogleLogin)
	}
	r.GET("/api/ouath/googleinfo/:token",controllers.SaveGoogleInfo)

	r.Run()
}
