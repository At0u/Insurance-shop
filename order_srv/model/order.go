package model

import "time"

type Order struct{
	BaseModel
	UserId int32 `gorm:"type:int;index;not null"`

	ProductID int32 `gorm:"type:int;not null" json:"productId"`
	ProductName string `gorm:"type:varchar(50);not null"`
	CategoryID int32 `gorm:"type:int;not null" json:"categoryId"`
	CategoryName string `gorm:"type:varchar(50);not null"`
	Term	int32 `gorm:"default:1;not null"`

	OrderSn string `gorm:"type:varchar(30);index"` //订单号，自己生成的订单号
	//status大家可以考虑使用iota来做
	Status int32 `gorm:"default:1;not null comment '1(履行中), 2(已过保期)， 3(提前解约), 4(状态异常)'"`
	Price float32	`gorm:"not null"`
	EndTime *time.Time `gorm:"type:datetime"`

	ApplicantName string `gorm:"type:varchar(50);not null" json:"applicantName"`
	ApplicantMobile string `gorm:"index:idx_applicant_mobile;type:varchar(11);not null" json:"applicantMobile"`
	ApplicantIdNum string `gorm:"index:idx_applicant_identity;varchar(50);not null" json:"applicantIdNum"`

	InsurerName string `gorm:"type:varchar(50);not null" json:"insurerName"`
	InsurerMobile string `gorm:"index:idx_insurer_mobile;type:varchar(11);not null" json:"insurerMobile"`
	InsurerIdNum string `gorm:"index:idx_insurer_identity;varchar(50);not null" json:"insurerIdNum"`
}

