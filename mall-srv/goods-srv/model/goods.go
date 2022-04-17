package model

type Category struct {
	BaseModel
	Name             string      `gorm:"varchar(20);not null" json:"name"`
	ParentCategoryId int32       `json:"parent"`
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryId;references:ID" json:"sub_category"`
	Level            int32       `gorm:"type:int;default:1;not null" json:"level"`
	IsTab            bool        `gorm:"default:false;not null" json:"is_tab"`
}

type Brand struct {
	BaseModel
	Name string `gorm:"varchar(20);not null"`
	Logo string `gorm:"varchar(200);default:'';not null"`
}

type GoodsCategoryBrands struct {
	BaseModel
	CategoryId int32 `gorm:"type:int;index:idx_category_brand;unique"`
	Category   Category
	BrandId    int32 `gorm:"type:int;index:idx_category_brand;unique"`
	Brand      Brand
}

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;default:1;not null"`
}

type Goods struct {
	BaseModel

	CategoryId int32 `gorm:"type:int;not null"`
	Category   Category
	BrandId    int32 `gorm:"type:int;not null"`
	Brand      Brand

	OnSale   bool `gorm:"default:false;not null"`
	ShipFree bool `gorm:"default:false;not null"`
	IsNew    bool `gorm:"default:false;not null"`
	IsHot    bool `gorm:"default:false;not null"`

	Name            string   `gorm:"type:varchar(50);not null"`
	GoodsSn         string   `gorm:"type:varchar(50);not null"`
	ClickNum        int32    `gorm:"type:int;default:0;not null"`
	SoldNum         int32    `gorm:"type:int;default:0;not null"`
	FavNum          int32    `gorm:"type:int;default:0;not null"`
	MarketPrice     float32  `gorm:"not null"`
	ShopPrice       float32  `gorm:"not null"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null"`
	Images          GormList `gorm:"type:varchar(1000);not null"`
	DescImages      GormList `gorm:"type:varchar(1000);not null"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null"`
	Stocks          int32    `gorm:"type:int;default:0;not null"`
}
