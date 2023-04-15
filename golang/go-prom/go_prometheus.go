package main

// go-prom refer to https://eddycjy.com/posts/prometheus/2020-05-16-metrics/
// it includes: metrics, counter, guage, histogram, summary
// and there is an demo about such 4 instance to describe the Go application.

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	prometheus "github.com/prometheus/client_golang/prometheus"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// gin.Engine to serve
	engi *gin.Engine

	// counter instance
	accessCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_request_total",
		},
		[]string{"method", "path"},
	)

	// gauge instance 仪表盘
	queueGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "queue_num_total",
		},
		[]string{"name"},
	)

	// histogram 直方图
	durationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_durations_histogram_seconds",   // 名称
			Buckets: []float64{0.2, 0.5, 1, 2, 5, 10, 30}, // 区间值
		},
		[]string{"path"},
	)

	// summary
	durationSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_durations_summary_seconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"path"},
	)
)

func init() {
	rand.Seed(time.Now().UnixNano())

	// register counter into metrics
	prometheus.MustRegister(accessCounter)
	prometheus.MustRegister(queueGauge)
	prometheus.MustRegister(durationHistogram)
	prometheus.MustRegister(durationSummary)

	// initlize engine .
	engi = gin.Default()

	// echo handler .
	engi.GET("/echo", func(c *gin.Context) {
		param := c.Query("p")
		c.String(http.StatusOK, param)
	})

	// prometheus handlers
	// counter
	engi.GET("/counter", func(c *gin.Context) {
		u, _ := url.Parse(c.Request.RequestURI)
		accessCounter.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   u.Path,
		}).Add(1)

		c.String(http.StatusOK, "counter")
	})
	// gauge
	engi.GET("/queue-guage", func(c *gin.Context) {
		num := c.Query("num")
		fnum, _ := strconv.ParseFloat(num, 64)
		queueGauge.With(prometheus.Labels{
			"name": "queue_guage",
		}).Set(fnum)

		c.String(http.StatusOK, "queue_guage")
	})
	// histogram
	engi.GET("/histogram", func(c *gin.Context) {
		u, _ := url.Parse(c.Request.RequestURI)
		durationHistogram.With(prometheus.Labels{
			"path": u.Path,
		}).Observe(rand.Float64())

		c.String(http.StatusOK, "histogram")
	})
	// summary .
	engi.GET("/summary", func(c *gin.Context) {
		u, _ := url.Parse(c.Request.RequestURI)
		f := rand.Float64()
		println("this is f=", f)
		durationSummary.With(prometheus.Labels{
			"path": u.Path,
		}).Observe(f)

		c.String(http.StatusOK, "summary")
	})

	// promtheus metrics
	engi.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func main() {
	if err := engi.Run(":8080"); err != nil {
		log.Fatalf("could not run on :8080, err=%v", err)
	}
}
