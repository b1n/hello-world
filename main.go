package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] - %s %s path=%s, status_code=%d, latency=%s, user_agent=%s, error_message=%s\n",
			param.TimeStamp.Format(time.RFC1123),
			param.ClientIP,
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.GET("/", func(c *gin.Context) {
		text := fmt.Sprintf("%s - %s",os.Getenv("TEST_TEXT"), uuid.New().String())
		c.String(http.StatusOK, text)
	})

	log.Printf("start web server on %s", os.Getenv("PORT"))

	if err := router.Run(":" + os.Getenv("PORT")); err != nil {
		log.Println("router run :"+os.Getenv("PORT"), err)
	}
}
