package model

import (
	"time"
)

type BaseModel struct {
	ID        int32     `gorm:"primaryKey;type:int"`
	CreatedAt time.Time `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"-"`
	// DeletedAt gorm.DeletedAt `json:"-"`
	IsDeleted bool `json:"-"`
}
