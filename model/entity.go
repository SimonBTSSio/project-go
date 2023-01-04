package model

import "time"

type Product struct {
	ID        int       `json:"id"`
	Name      string    `gorm:"unique"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Payment struct {
	ID        int `json:"id"`
	ProductId int
	Product   Product
	PricePaid float64   `json:"price_paid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
