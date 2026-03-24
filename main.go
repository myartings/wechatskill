package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.String("port", "8090", "Server port")
	flag.Parse()

	app := NewAppServer(*port)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	app.setupRoutes(r)

	addr := fmt.Sprintf(":%s", *port)
	log.Printf("WeChat MCP Server starting on %s", addr)
	log.Printf("MCP endpoint: http://localhost:%s/mcp", *port)
	log.Printf("REST API: http://localhost:%s/api/v1/", *port)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
