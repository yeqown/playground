package v2

import (
	"context"

	"github.com/playground/golang/gormcs"

	"github.com/opentracing/opentracing-go"
	gorm2 "gorm.io/gorm"
)

// https://gorm.io/zh_CN/docs/prometheus.html
// https://github.com/go-gorm/prometheus

type opentracingPlugin struct{}

func (o opentracingPlugin) AfterCreate(db *gorm2.DB) {
	sp := db.Statement.Context.Value(gormContextSpan).(opentracing.Span)
	sp.Finish()

	println("called after")

	return
}

func (o opentracingPlugin) BeforeCreate(db *gorm2.DB) {
	sp := opentracing.SpanFromContext(db.Statement.Context)
	if sp != nil {
		sp = opentracing.StartSpan(
			"gorm2/create",
			opentracing.ChildOf(sp.Context()),
		)
	} else {
		sp = opentracing.StartSpan("gorm2/create")
	}

	println("called before")
	db.Statement.Context = context.WithValue(db.Statement.Context, gormContextSpan, sp)

	return
}

func (o opentracingPlugin) Name() string {
	return "opentracing"
}

func (o opentracingPlugin) Initialize(db *gorm2.DB) error {
	println("opentracing initialized")

	// register callbacks
	_ = db.Callback().Create().Before("gorm:create").Register("opentracing:before_create", o.BeforeCreate)
	_ = db.Callback().Create().After("gorm:create").Register("opentracing:after_create", o.AfterCreate)

	// set tracer
	return nil
}

var gormContextSpan struct{}

func newPlugin() gorm2.Plugin {
	return opentracingPlugin{}
}

func (g gorm2TestSuite) Test_withPlugin() {
	err := g.db().Use(newPlugin())
	g.Require().Nil(err)
	g.Require().NotNil(g.db())

	loc := &gormcs.LocationModel{
		UserID:   1,
		Country:  "CN",
		Province: "Sichuan",
		City:     "Chengdu",
	}

	// create
	err = g.db().WithContext(context.Background()).Create(loc).Error
	g.Assert().Nil(err)
}
