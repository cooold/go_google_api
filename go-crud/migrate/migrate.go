package main

import (
	"github.com/Qbom/go-crud/initializers"
	"github.com/Qbom/go-crud/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

}

func main() {
	initializers.DB.AutoMigrate(&models.Post{})
}
