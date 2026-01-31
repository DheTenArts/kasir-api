package models

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
	CategoryID int `json:"-"`
	Category Category `json:"category"`
}

type Category struct {
	ID int `json:"id"`
	Name string `json:"name"`
}