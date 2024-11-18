package objects

import (
	"strconv"
	"strings"
)

/* INTERNAL REPRESENTATION */
type Receipt struct {
	Id             string // uuid string
	Retailer       string
	PurchaseYear   int16
	PurchaseMonth  int8
	PurchaseDay    int8
	PurchaseHour   int8
	PurchaseMinute int8
	Items          []Item
	Total          float64
	Points         int64
	TotalDollars   int64
	TotalCents     int64
}

type Item struct {
	ShortDescription string
	Price            float64
}

/* API LAYER REPRESENTATION */
type APIReceipt struct {
	Retailer     string    `json:"retailer"`
	PurchaseDate string    `json:"purchaseDate"` // (YYYY-MM-DD)
	PurchaseTime string    `json:"purchaseTime"` // (HH:MM)
	Items        []APIItem `json:"items"`
	Total        string    `json:"total"`
}

type APIItem struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

/*
Converts the API Layer receipt representation received in requests
to an internal representation of receipts, which is processed and stored.
*/
func Convert(api *APIReceipt) (Receipt, error) {
	var r Receipt

	r.Retailer = api.Retailer

	dateParts := strings.Split(api.PurchaseDate, "-")

	year, err := strconv.ParseInt(dateParts[0], 10, 16)
	if err != nil {
		return r, err
	}
	r.PurchaseYear = int16(year)

	month, err := strconv.ParseInt(dateParts[1], 10, 8)
	if err != nil {
		return r, err
	}
	r.PurchaseMonth = int8(month)

	day, err := strconv.ParseInt(dateParts[2], 10, 8)
	if err != nil {
		return r, err
	}
	r.PurchaseDay = int8(day)

	timeParts := strings.Split(api.PurchaseTime, ":")

	hour, err := strconv.ParseInt(timeParts[0], 10, 8)
	if err != nil {
		return r, err
	}
	r.PurchaseHour = int8(hour)

	minute, err := strconv.ParseInt(timeParts[1], 10, 8)
	if err != nil {
		return r, err
	}
	r.PurchaseMinute = int8(minute)

	r.Total, err = strconv.ParseFloat(api.Total, 64)
	if err != nil {
		return r, err
	}

	totalParts := strings.Split(api.Total, ".")
	r.TotalDollars, err = strconv.ParseInt(totalParts[0], 10, 64)
	if err != nil {
		return r, err
	}

	r.TotalCents, err = strconv.ParseInt(totalParts[1], 10, 64)
	if err != nil {
		return r, err
	}

	r.Items = make([]Item, len(api.Items))

	for i, element := range api.Items {
		r.Items[i] = Item{}
		r.Items[i].ShortDescription = element.ShortDescription
		r.Items[i].Price, err = strconv.ParseFloat(element.Price, 64)

		if err != nil {
			return r, err
		}
	}

	return r, nil
}
