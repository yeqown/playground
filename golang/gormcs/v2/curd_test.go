package v2

import (
	"testing"

	"github.com/playground/gonic/gormcs"

	"github.com/stretchr/testify/suite"
	gorm2 "gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TestRunAll(t *testing.T) {
	suite.Run(t, new(gorm2TestSuite))
}

func (g *gorm2TestSuite) Test_curd_simple() {
	loc := &gormcs.LocationModel{
		UserID:   1,
		Country:  "CN",
		Province: "Sichuan",
		City:     "Chengdu",
	}

	// create and query
	err := g.db().Create(loc).Error
	g.Assert().Nil(err)
	out := new(gormcs.LocationModel)
	err = g.db().First(out, loc.ID).Error
	g.Assert().Equal(loc.UserID, out.UserID)
	g.Assert().Equal(loc.Country, out.Country)
	g.Assert().Equal(loc.City, out.City)
	g.Assert().Equal(loc.Province, out.Province)

	// update and query
	err = g.db().Where("id = ?", loc.ID).
		Updates(gormcs.LocationModel{
			Province: "SiChuan",
			City:     "ChengDu",
		}).Error
	g.Assert().Nil(err)
	out = new(gormcs.LocationModel)
	err = g.db().Model(loc).First(out, loc.ID).Error
	g.Assert().Equal(loc.UserID, out.UserID)
	g.Assert().Equal(loc.Country, out.Country)
	g.Assert().Equal("ChengDu", out.City)
	g.Assert().Equal("SiChuan", out.Province)

	// delete and query
	err = g.db().Where("id = ?", loc.ID).Delete(out).Error
	g.Assert().Nil(err)
	out = new(gormcs.LocationModel)
	err = g.db().Model(loc).First(out, loc.ID).Error
	g.Assert().Equal(gorm2.ErrRecordNotFound, err)
}

// 使用 Preload 和 Joins 来关联查询数据
func (g *gorm2TestSuite) Test_curd_association() {
	user := gormcs.UserModel{
		Name: "association",
		Sex:  1,
		Location: gormcs.LocationModel{
			Country:  "CN",
			Province: "SC",
			City:     "CD",
		},
		Careers: []gormcs.CareerModel{
			{
				Syear: 1993,
				Eyear: 1994,
				Desc:  "The First Job",
			},
			{
				Syear: 1995,
				Eyear: 2000,
				Desc:  "The Second Job",
			},
		},
	}

	err := g.db().Create(&user).Error
	g.Assert().Nil(err)

	// Preload 适用于 1v1 1vN
	out := new(gormcs.UserModel)
	err = g.db().Model(&gormcs.UserModel{}).
		Preload("Location").
		Preload("Careers").
		First(out, user.ID).Error
	g.Assert().Nil(err)
	g.Assert().Equal(user.Location.Province, out.Location.Province)
	g.Assert().Equal(len(user.Careers), len(out.Careers))
	g.Assert().Equal(2, len(out.Careers))
	g.T().Logf("%+v", out)

	// Join Preload 适用于 1 v 1
	out = new(gormcs.UserModel)
	err = g.db().Model(&gormcs.UserModel{}).
		Joins("Location").
		Preload("Careers").
		First(out, user.ID).Error
	g.Assert().Nil(err)
	g.Assert().Equal(user.Location.Province, out.Location.Province)
	g.Assert().Equal(len(user.Careers), len(out.Careers))
	g.Assert().Equal(2, len(out.Careers))
	g.T().Logf("%+v", out)
}

func (g *gorm2TestSuite) Test_batchCreate() {
	locations := []gormcs.LocationModel{
		{
			UserID:   1,
			Country:  "CN",
			Province: "SC",
			City:     "CD",
		},
		{
			UserID:   1,
			Country:  "CN",
			Province: "SC",
			City:     "CD",
		},
		{
			UserID:   1,
			Country:  "CN",
			Province: "SC",
			City:     "CD",
		},
		{
			UserID:   1,
			Country:  "CN",
			Province: "SC",
			City:     "CD",
		},
	}

	// 批量插入
	err := g.db().CreateInBatches(locations, len(locations)).Error
	g.Assert().Nil(err)
}

func (g *gorm2TestSuite) Test_upsert() {
	g.T().Skip("should test on MySQL")

	career := &gormcs.CareerModel{
		UserID: 9999999,
		Syear:  2001,
		Eyear:  2005,
		Desc:   "No Job",
	}

	// upsert once
	err := g.db().Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "syear"},
			{Name: "eyear"},
		},
		Where: clause.Where{
			Exprs: []clause.Expression{
				clause.Eq{
					Column: "user_id",
					Value:  9999999,
				},
			},
		}, // UserId conflict
		//DoNothing: false,
		//DoUpdates: clause.AssignmentColumns([]),
		//UpdateAll: false,
	}).Create(career).Error
	g.Assert().Nil(err)

	// update and upsert again
	career.Syear = 2020
	career.Eyear = 2025
	career.Desc = "I'm Updated"
	err = g.db().Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "syear"},
			{Name: "eyear"},
		},
		Where: clause.Where{
			Exprs: []clause.Expression{
				clause.Eq{
					Column: "user_id",
					Value:  9999999,
				},
			},
		},
		//DoNothing: false,
		//DoUpdates: clause.AssignmentColumns([]),
		//UpdateAll: false,
	}).Create(career).Error
	g.Assert().Nil(err)

	// count
	count := int64(0)
	err = g.db().Model(career).Count(&count).Error
	g.Assert().Nil(err)
	g.Assert().Equal(1, count)

	// query and compare
	out := new(gormcs.CareerModel)
	err = g.db().Model(career).First(out, career.ID).Error
	g.Assert().Nil(err)
	g.Assert().Equal(career.Syear, out.Syear)
	g.Assert().Equal(career.Eyear, out.Eyear)
	g.Assert().NotEqual(career.Desc, out.Desc)
	g.Assert().NotEqual("No Job", out.Desc)

}

func (g *gorm2TestSuite) Test_scopesQuery() {
	withUserId := func(userId int64) func(db *gorm2.DB) *gorm2.DB {
		return func(db *gorm2.DB) *gorm2.DB {
			return db.Where("user_id = ?", userId)
		}
	}

	inCityCD := func(db *gorm2.DB) *gorm2.DB {
		return db.Where("city = ?", "CD")
	}

	out := new(gormcs.LocationModel)
	stm := g.db().
		Session(&gorm2.Session{DryRun: true}).
		Model(out).
		Scopes(
			withUserId(10),
			inCityCD,
		).
		First(out).Statement

	g.Assert().Equal(
		"SELECT * FROM `location` WHERE user_id = ? AND city = ? ORDER BY `location`.`id` LIMIT 1",
		stm.SQL.String(),
	)
	g.Assert().Equal(int64(10), stm.Vars[0])
	g.Assert().Equal("CD", stm.Vars[1])

	// A smart way to use scopes, should be define a builder to build
	// builder := NewBuilder(db)
	// builder.WithUserId().WithName().AgeGEQ(30)... and so on
	// db.WithScopes(builder.Scopes).Find()

}
