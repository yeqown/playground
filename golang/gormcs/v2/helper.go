package v2

import (
	"os"
	"time"

	"github.com/playground/gonic/gormcs"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	gorm2 "gorm.io/gorm"
)

func connectDB() (_db *gorm2.DB) {
	var err error

	// v2 配置
	config := &gorm2.Config{
		SkipDefaultTransaction:                   false,
		NamingStrategy:                           nil,
		FullSaveAssociations:                     false,
		Logger:                                   nil,
		NowFunc:                                  nil,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		CreateBatchSize:                          0,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	}

	// v2 链接到数据库
	_db, err = gorm2.Open(sqlite.Open("./testdata/sqlite3.db"), config)
	if err != nil {
		panic("failed to connect database")
	}

	// debug 开启
	_db = _db.Debug()
	sqlDB, err := _db.DB()
	if err != nil {
		panic("failed to get sql.DB")
	}

	// 链接池配置
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Minute)

	return _db
}

// prepare testdata in _db
func prepareTestdata(_db *gorm2.DB) {
	// clear up
	//_db.Drop(
	//	&gormcs.LocationModel{},
	//	&gormcs.UserModel{},
	//	&gormcs.CareerModel{},
	//)

	// migrate tables
	_ = _db.AutoMigrate(
		&gormcs.UserModel{},
		&gormcs.LocationModel{},
		&gormcs.CareerModel{},
	)
}

type gorm2TestSuite struct {
	suite.Suite

	_db *gorm2.DB
}

func (g *gorm2TestSuite) db() *gorm2.DB {
	return g._db
}

func (g *gorm2TestSuite) TearDownSuite() {
	err := os.Remove("./testdata/sqlite3.db")
	g.Assert().Nil(err)

}

func (g *gorm2TestSuite) SetupSuite() {
	g._db = connectDB()
	prepareTestdata(g._db)
}
