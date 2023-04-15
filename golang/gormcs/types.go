package gormcs

import "github.com/jinzhu/gorm"

// UserModel .
type UserModel struct {
	gorm.Model
	Name     string
	Sex      uint
	Location LocationModel `gorm:"foreignkey:user_id;"`
	Careers  []CareerModel `gorm:"foreignkey:user_id"`
}

// TableName of UserModel in DB
func (m UserModel) TableName() string {
	return "user"
}

// LocationModel .
type LocationModel struct {
	gorm.Model
	UserID   uint
	Country  string
	Province string
	City     string
}

// TableName of LocationModel in DB
func (m LocationModel) TableName() string {
	return "location"
}

// CareerModel .
type CareerModel struct {
	gorm.Model
	UserID uint
	Syear  uint
	Eyear  uint
	Desc   string
}

// TableName of CareerModel in DB
func (m CareerModel) TableName() string {
	return "career"
}
