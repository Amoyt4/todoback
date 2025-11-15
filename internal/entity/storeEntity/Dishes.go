package storeEntity

type Dish struct {
	ID     uint   `json:"id" example:"1"`
	Name   string `json:"name" example:"some name"`
	Price  int    `json:"price" example:"100"`
	ImgUrl string `json:"img_url"`
}

type NewDish struct {
	Name   string `json:"name" example:"some name"`
	Price  int    `json:"price" example:"100"`
	ImgUrl string `json:"img_url"`
}

type DeleteDish struct {
	ID uint `json:"id" example:"1"`
}
