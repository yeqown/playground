package v2

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/playground/golang/gormcs"
	gorm2 "gorm.io/gorm"

	"gorm.io/gorm/logger"
)

var (
	_loggerContextKey = struct{}{}
)

type tracingLoggerForGorm2 struct {
	std *log.Logger
}

func newTracingLogger(w io.Writer) tracingLoggerForGorm2 {
	if w == nil {
		w = os.Stdout
	}
	return tracingLoggerForGorm2{
		std: log.New(w, "tracing", log.Llongfile|log.LUTC),
	}
}

func (t tracingLoggerForGorm2) LogMode(level logger.LogLevel) logger.Interface {
	return t
}

func (t tracingLoggerForGorm2) Info(ctx context.Context, s string, i ...interface{}) {
	traceId := parseContextAsExample(ctx)
	log.Printf(s+traceId, i...)
}

func (t tracingLoggerForGorm2) Warn(ctx context.Context, s string, i ...interface{}) {
	traceId := parseContextAsExample(ctx)
	log.Printf(s+traceId, i...)
}

func (t tracingLoggerForGorm2) Error(ctx context.Context, s string, i ...interface{}) {
	traceId := parseContextAsExample(ctx)
	log.Printf(s+traceId, i...)
}

func (t tracingLoggerForGorm2) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	traceId := parseContextAsExample(ctx)
	s, i := fc()
	log.Printf("begin=%s, fc=%s, %d, err=%v traceId=%s", begin.String(), s, i, err, traceId)
}

func parseContextAsExample(ctx context.Context) string {
	v := ctx.Value(_loggerContextKey)
	return v.(string)
}

func (g gorm2TestSuite) Test_tracingLogger() {
	loc := &gormcs.LocationModel{
		UserID:   1111,
		Country:  "222",
		Province: "333",
		City:     "444",
	}

	// mock opentracing spanContext
	tracingCtx := context.WithValue(context.TODO(), _loggerContextKey, "FqKXVdTvIx_mPjOYdjDyUSy_H1jr")

	// do an operation
	buf := bytes.NewBuffer(nil)
	err := g.db().
		Session(&gorm2.Session{Logger: newTracingLogger(buf)}).
		WithContext(tracingCtx).
		Create(loc).Error
	g.Assert().Nil(err)
	g.T().Log(buf.String())
}
