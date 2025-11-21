package storeEntity

type OrderItem struct {
	ID       int `json:"id"`
	OrderID  int `json:"order_id"`
	DishID   int `json:"dish_id"`
	Quantity int `json:"quantity"`
}

type OrderItemCreate struct {
	DishID   int `json:"dish_id" validate:"required,min=1"`
	Quantity int `json:"quantity" validate:"required,min=1"`
}

type GetOrderItemsRequest struct {
	OrderID int `json:"order_id"`
}

type UpdateOrderItemRequest struct {
	OrderItemID int `json:"order_item_id"`
	Quantity    int `json:"quantity"`
}

type DeleteOrderItemRequest struct {
	OrderItemID int `json:"order_item_id"`
}
