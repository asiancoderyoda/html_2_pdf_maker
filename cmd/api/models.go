package main

import "encoding/json"

type Invoice struct {
	Id              int64         `json:"id"`
	InvoiceId       string        `json:"invoice_id"`
	Date            string        `json:"date"`
	CustomerId      string        `json:"customer_id"`
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

func (inv *Invoice) GetID() int64 {
	return inv.Id
}

func (item *InvoiceItem) GetID() int64 {
	return item.Id
}

type TemplateInterface interface {
	GetID() int64
}

type Request struct {
	Data json.RawMessage
}
