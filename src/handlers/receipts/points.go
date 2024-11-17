package receipts

import "github.com/gin-gonic/gin"

func GetPoints(c *gin.Context) {
	c.JSON(200, gin.H{"points": 32})
}
