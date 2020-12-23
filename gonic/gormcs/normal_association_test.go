package gormcs

import (
	"fmt"
	"testing"
)

// getUsersNormally .
func getUsersNormally() ([]UserModel, error) {
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

// saveUserNomally .
func saveUserNomally() error {
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

func TestGetUsersNormally(t *testing.T) {
	prepareTestdata()

	if us, err := getUsersNormally(); err != nil {
		fmt.Println(err)
	} else {
		// fmt.Printf("%v\n", us)
		for _, v := range us {
			if v.Location.ID == 0 {
				t.Errorf("could not association to Location")
			}

			if len(v.Careers) == 0 || v.Careers[0].ID == 0 {
				t.Errorf("could not association to Careers")
			}
		}
	}
}

func Test_normalCreate(t *testing.T) {
	prepareTestdata()

	if err := saveUserNomally(); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
