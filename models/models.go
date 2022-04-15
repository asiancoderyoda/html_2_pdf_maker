package models

type Invoice struct {
	Id              int64         `json:"id"`
	Number          string        `json:"number"`
	Date            string        `json:"date"`
	CustomerId      int64         `json:"customer_id"`
	CustomerName    string        `json:"customer_name"`
	CustomerAddress string        `json:"customer_address"`
	Total           float64       `json:"total"`
	Items           []InvoiceItem `json:"items"`
}

type InvoiceItem struct {
	Id          int64   `json:"id"`
	ItemName    string  `json:"item_name"`
	Description string  `json:"description"`
	Quantity    int64   `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}
