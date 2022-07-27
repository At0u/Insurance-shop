package model

type Category struct{
	BaseModel
	Name  string `gorm:"type:varchar(20);not null" json:"name"`
	IsBan int`gorm:"type:int;default:2;not null comment'1表示在商品页面分类中，2表示不在'"`
}

type Banner struct {
	BaseModel
	//图片
	Image string `gorm:"type:varchar(200);not null"`
	//跳转到对应商品的位置
	Url string `gorm:"type:varchar(200);not null"`
	Index int32 `gorm:"type:int;default:1;not null"`
}

type Goods struct {
	BaseModel

	AgeType int32 `gorm:"type:int;not null comment'1表示儿童，2表示成人,3是老人'"`
	Term	int32 `gorm:"default:1;not null"`
	CategoryID int32 `gorm:"type:int;not null"`
	Category Category

	OnSale bool `gorm:"default:false;not null"`
	IsNew bool `gorm:"default:false;not null"`
	IsHot bool `gorm:"default:false;not null"`

	Name  string `gorm:"type:varchar(50);not null"`
	GoodsSn string `gorm:"type:varchar(50);not null"`

	Price float32 `gorm:"not null"`
	GoodsBrief string `gorm:"type:varchar(100);not null"`
	Images GormList `gorm:"type:varchar(1000);not null"`
	DescImages GormList `gorm:"type:varchar(1000);not null"`
	GoodsFrontImage string `gorm:"type:varchar(200);not null"`
}


