// models.article.go

package main

import (
	"math/rand"
)

type application struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Status  string `json:"status"`
	Loan    int    `json:"loan"`
}

func createNewApplication(name, address string, loan int) *application {
	// Set the ID of a new article to one more than the number of articles
	a := application{ID: rand.Int(), Name: name, Address: address, Loan: loan, Status: "Pending"}
	return &a
}
