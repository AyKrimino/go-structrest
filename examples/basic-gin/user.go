package main

type User struct {
	ID    int    `json:"id" crud:"pk"`
	Name  string `json:"name" crud:"searchable"`
	Email string `json:"email" crud:"searchable"`
}
