package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type OrderItem struct {
	ItemID     int     `json:"item_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}

type UpdateOrderStatus struct {
	OrderTime string `json:"order_time"`
}

type UserId struct {
	ID int `json:"user_id"`
}
