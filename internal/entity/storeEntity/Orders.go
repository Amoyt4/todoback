package storeEntity

import "time"

type Order struct {
	ID            int         `json:"id"`
	RoomNum       int         `json:"room_num"`
	TimeToDeliver time.Time   `json:"time_to_deliver"`
	TotalSum      int         `json:"total_summ"`
	Items         []OrderItem `json:"items"`
}

type CreateOrderRequest struct {
	RoomNum       int               `json:"room_num" validate:"required,min=1"`
	TimeToDeliver time.Time         `json:"time_to_deliver" validate:"required"`
	Items         []OrderItemCreate `json:"items" validate:"required,min=1"`
}

type GetOrderByIDRequest struct {
	ID int `json:"id"`
}

type DeleteOrder struct {
	ID int `json:"id"`
}
