package v2

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	gorm2 "gorm.io/gorm"
)

type customTypeIds []uint32

// FIXME: could not be saved into database
func (c *customTypeIds) Value() (driver.Value, error) {
	byts, err := json.Marshal(c)
	return string(byts), err
}

func (c *customTypeIds) Scan(src interface{}) error {
	byts, ok := src.([]byte)
	if !ok {
		return errors.New("invalid type")
	}

	err := json.Unmarshal(byts, c)
	fmt.Printf("%+v", c)
	return err
}

type withCustomTypeModel struct {
	gorm2.Model

	Ids customTypeIds `gorm:"column:ids"`
}

func (w withCustomTypeModel) TableName() string {
	return "with_custom_type"
}

func (g gorm2TestSuite) Test_CustomType() {
	_ = g.db().AutoMigrate(&withCustomTypeModel{})

	// create
	m := &withCustomTypeModel{
		Ids: customTypeIds([]uint32{1, 2, 3, 4}),
	}
	err := g.db().Model(m).Create(m).Error
	g.Assert().Nil(err)

	// query
	out := new(withCustomTypeModel)
	err = g.db().Model(out).First(out, m.ID).Error
	g.Assert().Nil(err)
	g.Assert().Equal(m.Ids, out.Ids)
	g.Equal(4, len(out.Ids))
	g.T().Log(out.Ids)
}
