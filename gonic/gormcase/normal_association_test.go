package gormcase

import (
	"fmt"
	"testing"
)

func TestGetUsersNormally(t *testing.T) {
	if us, err := GetUsersNormally(); err != nil {
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
	if err := SaveUserNomally(); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
