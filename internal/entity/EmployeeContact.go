package entity

type EmployeeContact struct {
	Id      uint   `json:"id"`
	RoomNum uint   `json:"room_num"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
}

type NewEmployeeContact struct {
	RoomNum uint   `json:"room_num"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
}

type DeleteEmployeeContact struct {
	Id uint `json:"id"`
}
