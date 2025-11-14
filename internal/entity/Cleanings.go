package entity

import "time"

type Cleanings struct {
	ID        uint      `json:"id" example:"1"`
	RoomNum   int       `json:"room_num" example:"1"`
	StartTime time.Time `json:"start_time" example:"2020-01-01T00:00:00+00:00"`
	EndTime   time.Time `json:"end_time" example:"2020-01-01T00:00:00+00:00"`
	Comment   string    `json:"comment" example:"comment"`
}

type NewClean struct {
	RoomNum   int       `json:"room_num" example:"1"`
	StartTime time.Time `json:"start_time" example:"2020-01-01T00:00:00+00:00"`
	EndTime   time.Time `json:"end_time" example:"2020-01-01T00:00:00+00:00"`
	Comment   string    `json:"comment" example:"comment"`
}

type CleanDel struct {
	ID uint `json:"id" example:"1"`
}
