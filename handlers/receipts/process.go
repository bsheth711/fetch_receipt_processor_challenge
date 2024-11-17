package receipts

import (
	"service/objects"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var Receipts = make(map[string]objects.Receipt)

func ProcessReceipts(c *gin.Context) {

	var data objects.APIReceipt
	err := c.BindJSON(&data)

	if err != nil {
		c.Error(err)
		return
	}

	receipt, err := objects.Convert(&data)

	if err != nil {
		c.Error(err)
		return
	}

	receipt.Id = uuid.New().String()

	Receipts[receipt.Id] = receipt

	c.JSON(200, gin.H{"id": receipt.Id})
}

/*
func calculatePoints(c *objects.Receipt) {

}
*/
