package handler

import (
	pb "mall-srv/goods-srv/proto"

	"gorm.io/gorm"
)

type GoodsServer struct {
	pb.UnimplementedGoodsServer
}

// Paginate 分页逻辑
func Paginate(page, pageSize int32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}
