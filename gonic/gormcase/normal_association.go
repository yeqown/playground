package gormcase

import (
	"fmt"
)

// GetUsersNormally .
func GetUsersNormally() ([]UserModel, error) {
	var (
		us []UserModel
	)

	err := db.Model(&UserModel{}).Find(&us).Error
	for idx := range us {
		// find location
		if err := db.Model(&us[idx]).
			Related(&us[idx].Location, "Location").Error; err != nil {
			fmt.Println(err)
		}

		// find career
		if err := db.Model(&us[idx]).
			Related(&us[idx].Careers, "Careers").Error; err != nil {
			fmt.Println(err)
		}
	}
	return us, err
}

// SaveUserNomally .
func SaveUserNomally() error {
	user := &UserModel{
		Name: "yeqown",
		Sex:  1,
		Location: LocationModel{
			UserID:  1,
			Country: "China",
		},
		Careers: []CareerModel{
			{
				UserID: 1,
				Syear:  2019,
			},
		},
	}

	return db.Model(user).Create(user).Error
}
