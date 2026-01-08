## Proposed layout

.
├── dummy-app/
│   └── Dockerfile
├── monitor/
│   └── Dockerfile
├── prometheus/
│   └── prometheus.yml
├── grafana/
│   └── dashboards/
├── docker-compose.yml
└── README.md


## About

This project implements black-box monitoring using synthetic probes.
Availability SLO is computed using request success rate and latency thresholds.
Prometheus is used for metric aggregation and Grafana for visualization.


## TBD

The app has two endpoints /health and /work. The monitor written in go, that makes request on these two, and tell if the service is up, since how long, the latency, and number of non 200 responses and timeouts and serves it over /metrics. Monitor and app make use of goroutines. Both are containerised alongside prometheus and grafana to save the metrics and create dashboards depicting the metrics along with SLO that is availability( using latency, error rate and uptime). Broughtup with docker compose/docker, alerting in place to email or elsewhere. The monitor should shutdown cleanly.
