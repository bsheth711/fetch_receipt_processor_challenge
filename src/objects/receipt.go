package objects

import "strconv"

/* INTERNAL REPRESENTATION */
type Receipt struct {
	Id           string // uuid string representation
	Retailer     string
	PurchaseDate string // (YYYY-MM-DD)
	PurchaseTime string // (HH:MM)
	Items        []Item
	Total        float64
	Points       int64
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
Converts between the API Layer receipt representation received in requests
to an internal representation of receipts.
*/
func Convert(api *APIReceipt) (Receipt, error) {
	var r Receipt

	r.Retailer = api.Retailer
	r.PurchaseDate = api.PurchaseDate
	r.PurchaseTime = api.PurchaseTime

	total, err := strconv.ParseFloat(api.Total, 64)

	if err != nil {
		return r, err
	}

	r.Total = total

	r.Items = make([]Item, len(api.Items))

	for i, element := range api.Items {
		r.Items[i] = Item{}
		r.Items[i].ShortDescription = element.ShortDescription
		price, err := strconv.ParseFloat(element.Price, 64)

		if err != nil {
			return r, err
		}

		r.Items[i].Price = price
	}

	return r, err
}
