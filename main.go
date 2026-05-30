package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type sendRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func listenAddr() string {
	host := os.Getenv("BIND_ADDR")
	if host == "" {
		host = os.Getenv("HTTP_HOST")
	}
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("HTTP_PORT")
	}
	if port == "" {
		port = "8080"
	}

	return host + ":" + port
}

func main() {
	_ = godotenv.Load()

	r := gin.Default()

	r.POST("/mail/test", func(c *gin.Context) {
		var req sendRequest
		_ = c.ShouldBindJSON(&req)

		if err := sendTest(req.To, req.Subject, req.Body); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	if err := r.Run(listenAddr()); err != nil {
		panic(err)
	}
}
