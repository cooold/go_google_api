package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Qbom/go-crud/initializers"
	"github.com/Qbom/go-crud/models"
	"github.com/gin-gonic/gin"
)

func PostsCreate(c *gin.Context) {

	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	post := models.Post{Title: body.Title, Body: body.Body}
	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostsIndex(c *gin.Context) {
	var posts []models.Post
	initializers.DB.Find(&posts)
	c.JSON(200, gin.H{
		"posts": posts,
	})
}

func PostsShow(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	initializers.DB.First(&post, id)

	c.JSON(200, gin.H{
		"posts": post,
	})

}

func PostUpdate(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	var post models.Post
	initializers.DB.First(&post, id)

	initializers.DB.Model(&post).Updates(models.Post{
		Title: body.Title,
		Body:  body.Body,
	})

	c.JSON(200, gin.H{
		"posts": post,
	})

}

func PostDelet(c *gin.Context) {
	id := c.Param("id")

	initializers.DB.Delete(&models.Post{}, id)
	c.Status(200)
}

func TMDBList(c *gin.Context) {
	page := c.Param("page")
	url := "https://api.themoviedb.org/3/tv/top_rated?api_key=853e8506abed92c9e3f1c8c95e7afe8d&language=zh-TW&page=" + page
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	var data models.TMDBTopRated
	json.Unmarshal(body, &data)

	c.JSON(200, gin.H{
		"movie": data.Results,
	})
}

func SaveGoogleInfo(c *gin.Context) {
	token := c.Param("token")
	var name string
	var picture string
	var email string
	for i, part := range strings.Split(token, ".") {
		fmt.Printf("[%d] part: %s\n", i, part)
		decoded, err := base64.RawURLEncoding.DecodeString(part)
		if err != nil {
			panic(err)
		}
		fmt.Println("decoded:", string(decoded))
		if i != 1 {
			continue // i == 1 is the payload
		}

		var m map[string]interface{}
		if err := json.Unmarshal(decoded, &m); err != nil {
			fmt.Println("json decoding failed:", err)
			continue
		}
		if names, ok := m["name"]; ok {
			name = names.(string)
			fmt.Println("name:", names)
		}
		if pictures, ok := m["picture"]; ok {
			picture = pictures.(string)
			fmt.Println("name:", pictures)
		}
		if emails, ok := m["email"]; ok {
			email = emails.(string)
			fmt.Println("name:", emails)
		}
	}
	c.JSON(200, gin.H{
		"name": name,
		"picture": picture,
		"email": email,
	})
	

}
