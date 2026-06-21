package main

type User struct {
	ID    int    `json:"id" bun:",pk,autoincrement" crud:"pk,autoincrement"`
	Name  string `json:"name" bun:",notnull" crud:"searchable"`
	Email string `json:"email" bun:",notnull" crud:"searchable"`
}
