package gormcase

import (
	"fmt"
	"testing"
)

func TestGetUsersWithAssociation(t *testing.T) {
	if us, err := GetUsersWithPreload(); err != nil {
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
