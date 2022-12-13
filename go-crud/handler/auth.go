package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Qbom/go-crud/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/Qbom/go-crud/pkg/res"
)

const redirectURL = "http://localhost:3010/api/ouath/google/login"
const scope = "https://www.googleapis.com/auth/userinfo.profile"

// GoogleAccsess GoogleAccsess
func GoogleAccsess(c *gin.Context) {
	res.Success(c, gin.H{
		"url": oauthURL(),
	})
}

func oauthURL() string {
	u := "https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&response_type=code&scope=%s&redirect_uri=%s"

	return fmt.Sprintf(u, config.Val.GoogleClientID, scope, redirectURL)
}

// GoogleLogin GoogleLogin
func GoogleLogin(c *gin.Context) {
	code := c.Query("code")

	token, err := accessToken(code)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Debug("accessToken error")
		c.Redirect(http.StatusFound, "/")
		return
	}

	id, name, picture, err := getGoogleUserInfo(token)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Debug("getGoogleUserInfo error")
		c.Redirect(http.StatusFound, "/")
		return
	}
	c.JSON(200, gin.H{
		"id": id,
		"name": name,
		"picture": picture,
	})
	log.Infof("id: %v, name: %v, picture: %v", id, name,picture)
}

func accessToken(code string) (token string, err error) {
	u := "https://www.googleapis.com/oauth2/v4/token"

	data := url.Values{"code": {code}, "client_id": {config.Val.GoogleClientID}, "client_secret": {config.Val.GoogleSecretKey}, "grant_type": {"authorization_code"}, "redirect_uri": {redirectURL}}
	body := strings.NewReader(data.Encode())

	resp, err := http.Post(u, "application/x-www-form-urlencoded", body)
	if err != nil {
		return token, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return token, err
	}

	token = gjson.GetBytes(b, "access_token").String()

	return token, nil
}

func getGoogleUserInfo(token string) (id, name, picture string, err error) {
	u := fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=%s", token)
	resp, err := http.Get(u)
	if err != nil {
		return id, name, picture, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return id, name, picture, err
	}

	name = gjson.GetBytes(body, "name").String()
	id = gjson.GetBytes(body, "id").String()
	picture = gjson.GetBytes(body, "picture").String()
	log.Println(gjson.ParseBytes(body))
	return id, name, picture, nil
}
