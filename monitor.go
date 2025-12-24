package main

import (
	"flag"
	"net/http"
    "gopkg.in/yaml.v3"
	"time"
	"os"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var serviceUp = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "service_up",
		Help: "Service up status",
	},
	[]string{"app"},
)
var healthTimeout = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "health_timeout",
		Help: "health check timeout",
	},
	[]string{"app"},
)
var upSince = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "up_since",
		Help: "service uptime",
	},
	[]string{"app"},
)
var workLatency = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "app_latency",
		Help: "latency for the application",
	},
	[]string{"app"},
)
var failureCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "failure_count",
		Help: "app returning non 2xx",
	},
	[]string{"app"},
)
var appTimeout = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "app_timeout",
		Help: "application timeout",
	},
	[]string{"app"},
)
type Service struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}
type Monitor struct {
	Interval time.Duration `yaml:"interval"`
	Timeout time.Duration `yaml:"timeout"`
}
type Config struct {
	Services []Service `yaml:"services"`
	Monitor Monitor `yaml:"monitor"`
}

func main() {
	config := flag.String("config","/etc/monitor/config.yaml","config file location")
	flag.Parse()
	data, err := os.ReadFile(*config)
	if err != nil {
		panic(err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data,&cfg);err !=nil {
		panic(err)
	}
	prometheus.MustRegister(serviceUp)
	prometheus.MustRegister(healthTimeout)
	prometheus.MustRegister(upSince)
	prometheus.MustRegister(workLatency)
	prometheus.MustRegister(failureCount)
	prometheus.MustRegister(appTimeout)
	go func() {
		ticker := time.NewTicker(cfg.Monitor.Interval)
		defer ticker.Stop()
		for range ticker.C {
			for _, app := range cfg.Services {
				checkHealth(app.Name, app.URL, cfg.Monitor.Timeout)
				checkWork(app.Name, app.URL, cfg.Monitor.Timeout)
			}
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8081", nil)
}

func checkHealth(name string, url string, timeout time.Duration){

}

func checkWork(name string, url string, timeout time.Duration){

}