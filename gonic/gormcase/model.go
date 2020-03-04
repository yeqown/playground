package gormcase

import (
	"fmt"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func init() {
	connectDB()
	prepareTestdata()
}

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

var (
	db *gorm.DB
)

func connectDB() {
	var err error
	db, err = gorm.Open("sqlite3", "sqlite3.db")
	if err != nil {
		panic("failed to connect database")
	}
}

// prepare testdata in db
func prepareTestdata() {
	// clear up
	db.DropTableIfExists(&UserModel{}, &LocationModel{}, &CareerModel{})

	// migrate tables
	db.AutoMigrate(&UserModel{}, &LocationModel{}, &CareerModel{})

	user := &UserModel{
		Name: "yeqown",
		Sex:  1,
		Location: LocationModel{
			Country:  "CN",
			Province: "SC",
			City:     "Chengdu",
		},
		Careers: []CareerModel{
			{Syear: 13, Eyear: 14, Desc: "stage 1"},
			{Syear: 14, Eyear: 15, Desc: "stage 2"},
		},
	}

	// insert one user
	if err := db.Model(&UserModel{}).Create(user).Error; err != nil {
		fmt.Println(err)
	}

	// finish
	fmt.Println("prepareTestdata done")
}
