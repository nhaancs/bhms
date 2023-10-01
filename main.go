package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	server := http.Server{
		Handler: ginEngine,
		Addr:    ":3000",
	}

	const prefix = "api"
	router := ginEngine.Group(prefix)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, struct {
			Status string
		}{"OK"})
	})

	log.Println("Server start at port 3000")
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Server error: %+v\n", err)
		os.Exit(1)
	}
}
