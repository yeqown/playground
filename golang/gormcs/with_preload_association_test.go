package gormcs

import (
	"fmt"
	"testing"
)

// getUsersWithPreload .
func getUsersWithPreload() ([]UserModel, error) {
	var (
		us []UserModel
	)

	// find wanted users first
	db.Model(&UserModel{}).Find(&us)

	// association to another fields
	err := db.Preload("Location").Preload("Careers").Find(&us).Error
	return us, err
}

func TestGetUsersWithAssociation(t *testing.T) {
	prepareTestdata()

	if us, err := getUsersWithPreload(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", us)
		for _, v := range us {
			if v.Location.ID == 0 {
				t.Errorf("could not association to Location")
				t.Fail()
			}

			if len(v.Careers) == 0 || v.Careers[0].ID == 0 {
				t.Errorf("could not association to Careers")
				t.Fail()
			}
		}
	}
}
