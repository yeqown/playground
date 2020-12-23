package v2

import (
	"context"
	"errors"

	gorm2 "gorm.io/gorm"
)

var (
	_withContextKey = struct{}{}
)

type withContextTestModel struct {
	gorm2.Model

	Desc string `gorm:"column:desc"`
}

// ctx could be used in 'Callback / Hooks / Plugins / Logger'
func (w withContextTestModel) BeforeCreate(tx *gorm2.DB) error {
	ctx := tx.Statement.Context
	v := ctx.Value(_withContextKey)
	s := v.(string)
	if s == "" {
		return errors.New("ctx with empty s")
	}

	return nil
}

func (w withContextTestModel) TableName() string {
	return "with_context_test"
}

func (g gorm2TestSuite) Test_Context() {
	_ = g.db().AutoMigrate(&withContextTestModel{})

	root := context.Background()
	ctx := context.WithValue(root, _withContextKey, "test value")

	out := new(withContextTestModel)
	err := g.db().
		WithContext(ctx).
		Model(out).
		Create(out).Error
	g.Assert().Nil(err)
}
