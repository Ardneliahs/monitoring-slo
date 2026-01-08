package main

import (
	"flag"
	"net/http"
    "gopkg.in/yaml.v3"
	"time"
	"os"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"context"
	"encoding/json"
	"fmt"
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
var healthFailure = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "health_failure",
		Help: "health check failure",
	},
	[]string{"app"},
)
var appUnreachable = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "app_unreachable",
		Help: "monitor could not contact the app",
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
type HealthResponse struct {
	Status string `json:"status"`
	UptimeSec float64 `json:"uptime_sec"`
	Version string `json:"version"`
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
	registry := prometheus.NewRegistry()
	registry.MustRegister(serviceUp)
	registry.MustRegister(healthTimeout)
	registry.MustRegister(healthFailure)
	registry.MustRegister(appUnreachable)
	registry.MustRegister(upSince)
	registry.MustRegister(workLatency)
	registry.MustRegister(failureCount)
	registry.MustRegister(appTimeout)
	for _, svc := range cfg.Services {
		serviceUp.With(prometheus.Labels{"app": svc.Name}).Set(0)
		healthTimeout.With(prometheus.Labels{"app": svc.Name}).Set(0)
		healthFailure.With(prometheus.Labels{"app": svc.Name}).Set(0)
		appUnreachable.With(prometheus.Labels{"app": svc.Name}).Set(0)
		upSince.With(prometheus.Labels{"app": svc.Name}).Set(0)
		workLatency.With(prometheus.Labels{"app": svc.Name}).Set(0)
		failureCount.With(prometheus.Labels{"app": svc.Name}).Add(0)
		appTimeout.With(prometheus.Labels{"app": svc.Name}).Set(0)
    }
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
	handler := promhttp.HandlerFor(
    registry,
    promhttp.HandlerOpts{},
    )
	http.Handle("/metrics", handler)
	http.ListenAndServe(":8081", nil)
}

func checkHealth(name string, url string, timeout time.Duration){
	healthTimeout.With(prometheus.Labels{"app": name}).Set(0)
	healthFailure.With(prometheus.Labels{"app": name}).Set(0)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url+"/health",
		nil,
	)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			healthTimeout.WithLabelValues(name).Set(1)
			return
		}
		healthFailure.WithLabelValues(name).Set(1)
		return
	}
	defer resp.Body.Close()
	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		fmt.Println("JSON decode failed:", err)
		return
	}
	if health.Status == "UP" {
		serviceUp.WithLabelValues(name).Set(1)
	} else {
		serviceUp.WithLabelValues(name).Set(0)
	}
	upSince.WithLabelValues(name).Set(health.UptimeSec)
}

func checkWork(name string, url string, timeout time.Duration){
	appTimeout.With(prometheus.Labels{"app": name}).Set(0)
	appUnreachable.With(prometheus.Labels{"app": name}).Set(0)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url+"/work",
		nil,
	)
	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	latency := time.Since(start)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			appTimeout.WithLabelValues(name).Set(1)
			return
		}
		appUnreachable.WithLabelValues(name).Set(1)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		failureCount.WithLabelValues(name).Inc()
	} else {
		workLatency.WithLabelValues(name).Set(latency.Milliseconds() * 1000)
	}
}