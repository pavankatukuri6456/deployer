package main

import (
	"deployer/applier/server"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	server.Router(r)
	r.Run(":8080")
}
