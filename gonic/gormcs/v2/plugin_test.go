package v2

import gorm2 "gorm.io/gorm"

// https://gorm.io/zh_CN/docs/prometheus.html
// https://github.com/go-gorm/prometheus

type opentracingPlugin struct{}

func (o opentracingPlugin) Name() string {
	return "opentracing"
}

func (o opentracingPlugin) Initialize(db *gorm2.DB) error {
	panic("implement me")
}

func newPlugin() gorm2.Plugin {
	return opentracingPlugin{}
}

func (g gorm2TestSuite) Test_withPlugin() {
	err := g._db.Use(newPlugin())
	g.Require().Nil(err)
}
