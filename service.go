package main

import (
	"fmt"
	"service/handlers/receipts"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello World")

	r := gin.Default()
	r.POST("/receipts/process", receipts.ProcessReceipts)
	r.GET("/receipts/:id/points", receipts.GetPoints)

	r.Run(":8080")
}
