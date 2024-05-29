package main

import (
	"fmt"
	"html/template"
	"time"
)

type Customer struct {
	Name     string
	VatID    string
	Address1 string
	Address2 string
}

type Invoice struct {
	PlaceOfIssue string
	DateOfIssue  time.Time
	DateOfSale   time.Time
	Seller       Customer
	Buyer        Customer
	InvoiceNo    string
}

var seller = Customer{
	Name:     "ACME Corp.",
	VatID:    "0000000000",
	Address1: "42 Something Drive",
	Address2: "00000 New York City, NY",
}

var buyer = Customer{
	Name:     "XYZ Corp.",
	VatID:    "000000000",
	Address1: "42 Something Street",
	Address2: "Shibuya Ward, Tokyo",
}

type invoiceAssigns struct {
	Styles  template.CSS
	Invoice Invoice
}

func generateInvoice(date time.Time) Invoice {
	y, m, _ := date.Date()
	lastDay := time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC)

	invoiceNo := fmt.Sprintf("01/%02d/%d", m, y)

	return Invoice{
		PlaceOfIssue: "New York City",
		DateOfIssue:  lastDay,
		DateOfSale:   lastDay,
		Seller:       seller,
		Buyer:        buyer,
		InvoiceNo:    invoiceNo,
	}
}
