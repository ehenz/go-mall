package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type OrderDetailItem struct {
	GoodsId  int32
	GoodsNum int32
}

type OrderDetailList []OrderDetailItem

func (g *OrderDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

func (g OrderDetailList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

type BaseModel struct {
	ID        int32     `gorm:"primaryKey;type:int"`
	CreatedAt time.Time `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"-"`
	// DeletedAt gorm.DeletedAt `json:"-"`
	IsDeleted bool `json:"-"`
}
