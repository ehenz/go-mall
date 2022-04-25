package model

type Stock struct {
	BaseModel
	GoodsId int32 `gorm:"type:int;index"`
	Stock   int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"`
}

type OrderStatus struct {
	OrderSn string          `gorm:"type:varchar(200);index:order_sn,unique"`
	Status  int32           `gorm:"type:int"` // 1已扣减，2已归还
	Detail  OrderDetailList `gorm:"type:varchar(200)"`
}
