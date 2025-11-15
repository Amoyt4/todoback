package storeEntity

type OrderItem struct {
	ID       uint `json:"id" example:"1"`
	OrderID  uint `json:"order_id" example:"1"`
	DishID   uint `json:"dish_id" example:"1"`
	Quantity int  `json:"quantity" example:"2"`
}

type OrderItemWithDish struct {
	ID       uint `json:"id" example:"1"`
	OrderID  uint `json:"order_id" example:"1"`
	DishID   uint `json:"dish_id" example:"1"`
	Quantity int  `json:"quantity" example:"2"`
	Dish     Dish `json:"dish"`
}

type CreateOrderItem struct {
	OrderID  uint `json:"order_id" binding:"required" example:"1"`
	DishID   uint `json:"dish_id" binding:"required" example:"1"`
	Quantity int  `json:"quantity" binding:"required,min=1" example:"2"`
}

type UpdateOrderItem struct {
	Quantity int `json:"quantity" binding:"required,min=1" example:"3"`
}

type DeleteOrderItem struct {
	ID uint `json:"id" example:"1"`
}
