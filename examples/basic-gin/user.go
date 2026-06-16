package main

type User struct {
	ID    int    `json:"id" crud:"pk"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
