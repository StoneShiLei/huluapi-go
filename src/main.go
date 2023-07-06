package main

import (
	"huluapi/src/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("opencomputer", handler.OpenComputerHandler)
	r.POST("closecomputer", handler.CloseComputerHandler)

	r.Run("0.0.0.0:8080")
}
