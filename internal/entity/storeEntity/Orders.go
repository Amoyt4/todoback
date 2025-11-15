package storeEntity

import "time"

type Order struct {
	ID            uint                `json:"id" example:"1"`
	RoomNum       int                 `json:"room_num" example:"101"`
	TimeToDeliver time.Time           `json:"time_to_deliver"`
	TotalSum      int                 `json:"total_sum" example:"1500"`
	Items         []OrderItemWithDish `json:"items"`
}

type NewOrder struct {
	RoomNum       int       `json:"room_num" example:"1"`
	TimeToDeliver time.Time `json:"time_to_deliver"`
	TotalSum      int       `json:"total_sum" example:"1"`
}

type DeleteOrder struct {
	ID uint `json:"id" example:"1"`
}
