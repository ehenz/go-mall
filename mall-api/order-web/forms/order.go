package forms

type NewOrderForm struct {
	Name    string `json:"name" binding:"required,min=1,max=10"`
	Mobile  string `json:"mobile" binding:"required"`
	Address string `json:"address" binding:"required"`
	Post    string `json:"post" binding:"min=1,max=100"`
}
