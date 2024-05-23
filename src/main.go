package main

import (
	"huluapi/src/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("opencomputer", handler.OpenComputerHandler)
	r.POST("closecomputer", handler.CloseComputerHandler)
	r.POST("test", handler.Test)
	r.GET("simpleopencomputer", handler.SimpleOpenComputerHandler)

	r.Run("0.0.0.0:7096")
}
