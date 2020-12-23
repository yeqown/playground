package gormcs

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db   *gorm.DB
	once sync.Once
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
	once.Do(func() {
		connectDB()

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
	})
}
