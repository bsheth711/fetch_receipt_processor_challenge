package main

import (
	"service/src/handlers/receipts"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/receipts/process", receipts.ProcessReceipts)
	r.GET("/receipts/:id/points", receipts.GetPoints)

	r.Run(":8080")
}
