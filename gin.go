package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func setupRouter() *gin.Engine {
	gin.DisableConsoleColor()
	r := gin.New()
	return r
}

func getStatus(c *gin.Context) {
	msg := map[string]interface{}{"Status": "Ok", "msg": "ready", "version": "v1.0.1"}
	c.JSON(http.StatusOK, msg)
}

// func getCount(c *gin.Context) {
// 	playerstatsCache.ItemCount()
// 	msg := map[string]interface{}{"Status": "Ok", "guids": playerstatsCache.ItemCount()}
// 	c.JSON(http.StatusOK, msg)
// }

func main() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	r := setupRouter()
	log.Info("Starting Application")
	r.GET("/ready", getStatus)

	r.Run()
}
