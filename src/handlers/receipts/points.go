package receipts

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Handler for the endpoint /receipts/{id}/points.
Accepts a receipt id from /receipts/process's response.
Responds with the number of points associated with the receipt.
If no receipt with the specified id is found, responds with a 404 error.
*/
func GetPoints(c *gin.Context) {
	id := c.Param("id")
	receipt, ok := Receipts[id]

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Unrecognized receipt id."})
		c.Error(fmt.Errorf("unrecognized receipt id: %v", id))
		return
	}

	c.JSON(200, gin.H{"points": receipt.Points})
}
