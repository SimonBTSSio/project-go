package product

type InputProduct struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}

type InputProductUpdate struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
