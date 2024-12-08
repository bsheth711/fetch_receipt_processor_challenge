package receipts

import (
	"hash/fnv"
	"math"
	"net/http"
	"service/src/objects"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// in memory storage of receipt information
var Receipts = make(map[string]objects.Receipt)

// Stores hash for all received receipts.
// Used to prevent duplicate requests.
var UniqueReceipts = make(map[uint32]uint32)

/*
Handler for the endpoint /receipts/process.
Accepts a JSON payload in the format of APIReceipt.
Responds with an id for the newly created receipt,
or an error code if processing failed.
*/
func ProcessReceipts(c *gin.Context) {

	var data objects.APIReceipt
	err := c.BindJSON(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "Unable to deserialize JSON payload."})
		c.Error(err)
		return
	}

	hash := generateHash(&data)

	_, ok := UniqueReceipts[hash]

	if ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "Duplicate Request."})
		c.Error(err)
	}

	receipt, err := objects.Convert(&data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_SERVER_ERROR", "message": "Failed to process receipt."})
		c.Error(err)
		return
	}

	receipt.Id = uuid.New().String()
	receipt.Points = calculatePoints(&receipt)

	Receipts[receipt.Id] = receipt
	UniqueReceipts[hash] = hash

	c.JSON(http.StatusOK, gin.H{"id": receipt.Id})
}

/*
Returns int32 hash for a string.
*/
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

/*
Generates a unique hash based on an APIReceipt.
*/
func generateHash(receipt *objects.APIReceipt) uint32 {
	var str string

	str += receipt.Retailer
	str += receipt.PurchaseDate
	str += receipt.PurchaseTime
	str += receipt.Total

	for _, element := range receipt.Items {
		str += element.ShortDescription
		str += element.Price
	}

	return hash(str)
}

/*
Calculates points for a Receipt according to the following business logic:
1. One point for every alphanumeric character in the retailer name.
2. 50 points if the total is a round dollar amount with no cents.
3. 25 points if the total is a multiple of 0.25.
4. 5 points for every two items on the receipt.
5. If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
6. 6 points if the day in the purchase date is odd.
7. 10 points if the time of purchase is after 2:00pm and before 4:00pm.
8. 10 points if the character starts with the letter 'g' or 'G' only if the number of letters in the item description is not a multiple of 3
*/
func calculatePoints(receipt *objects.Receipt) int64 {
	var points int64 = 0

	// Rule 1
	for _, r := range receipt.Retailer {
		if isAlphanumeric(r) {
			points += 1
		}
	}

	// Rule 2
	if receipt.TotalCents == 0 {
		points += 50
	}

	// Rule 3
	if receipt.TotalCents%25 == 0 {
		points += 25
	}

	// Rule 4
	points += int64(len(receipt.Items) / 2 * 5)

	// Rule 5
	for _, element := range receipt.Items {
		trimmed := strings.TrimSpace(element.ShortDescription)

		if len(trimmed)%3 == 0 {
			points += int64(math.Ceil(element.Price * 0.2))
		} else {
			if trimmed[0] == 'G' || trimmed[0] == 'g' {
				points += 10
			}
		}
	}

	// Rule 6
	if receipt.PurchaseDay%2 == 1 {
		points += 6
	}

	// Rule 7
	after2PM := receipt.PurchaseHour > 14 || (receipt.PurchaseHour == 14 && receipt.PurchaseMinute > 0)
	before4PM := receipt.PurchaseHour < 16

	if after2PM && before4PM {
		points += 10
	}

	return points
}

/*
Helper function that takes a rune
and returns true if the rune is alphanumeric;
false otherwise.
*/
func isAlphanumeric(r rune) bool {
	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
		return true
	}

	return false
}
