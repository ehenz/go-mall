package forms

type CartItem struct {
	GoodsId int32 `json:"goods_id" binding:"required"`
	Nums    int32 `json:"nums" binding:"required,min=1"`
}

type CartUpdate struct {
	Nums    int32 `json:"nums" binding:"required,min=1"`
	Checked bool  `json:"checked"`
}
