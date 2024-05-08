package models

type MenuItem struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Available   bool    `json:"available"`
}

type Order struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id"`
	ItemID     int     `json:"item_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
	OrderTime  string  `json:"order_time"`
}
