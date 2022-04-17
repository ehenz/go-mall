package model

type Stock struct {
	BaseModel
	GoodsId int32 `gorm:"type:int;index"`
	Stock   int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"`
}
