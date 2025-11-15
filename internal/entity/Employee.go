package entity

type Employee struct {
	ID       uint   `json:"id" example:"1"`
	Login    string `json:"login" example:"admin"`
	Password string `json:"password" example:"123456"`
	Name     string `json:"name" example:"Employee"`
}

type NewEmployee struct {
	Login    string `json:"login" example:"admin"`
	Password string `json:"password" example:"123456"`
	Name     string `json:"name" example:"Employee"`
}

type DeleteEmployee struct {
	ID uint `json:"id" example:"1"`
}
