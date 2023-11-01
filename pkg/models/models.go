package models

import "time"

type Message struct {
	From      string    `json:"from"`
	To        []string  `json:"to"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type User struct {
	Email       string `json:"email"`
	Password    string `json:"-"`
	Position    string `json:"position"`
	Department  string `json:"department"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address,omitempty"`
	Salary      int    `json:"salary,omitempty"`
}
