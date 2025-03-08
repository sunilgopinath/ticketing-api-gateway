# ticketing-api-gateway

I'm going to use hello interview as a basis and implement all the things they talked about 

https://www.hellointerview.com/learn/system-design/problem-breakdowns/ticketmaster

## Getting started

```bash
$> git clone git@github.com:sunilgopinath/ticketing-api-gateway.git
$> cd ticketing-api-gateway
$> go run cmd/main.go -port=8080 -instance=gateway-1
{"level":"info","timestamp":"2025-03-04T19:47:54.903-0800","caller":"cmd/main.go:22","msg":"API Gateway is starting...","port":"8080","instance":"gateway-1"}
```

### Start jaeger (for distributed/microservice tracing)

```bash
$> docker run --rm -d --name jaeger \\n  -p 16686:16686 -p 4317:4317 -p 4318:4318 \\n  jaegertracing/all-in-one:latest
```

### Start REDIS (for cache)

```bash
$> docker run --rm -d --name redis -p 6379:6379 redis:latest
```

## How to use

### Start Application

```bash


```bash
$> curl "http://localhost:8080/bookings?user_id=125"
$> curl -X GET http://localhost:8080/events
$> curl -X POST http://localhost:8080/purchase
```

### Distributed tracing

[Open telemetry](https://opentelemetry.io/docs/languages/go/getting-started/) is used to collect traces for the microservices. Currently the collector is pointing to [JAEGER](https://www.jaegertracing.io/) for collection but in production we would point it at the open telemetry collector which would then sync with a grafana tempo.

The jaeger traces can be seen at http://localhost:16686/

### REDIS Cache

API requests are cached via REDIS. The API-gateway calls the appropriate handler and the handler checks the REDIS cache for the result. The cache key is an encoding of the URL parameters with an endpoint prefix

### Metrics/Dashboard

Ensure clean start

```bash
docker stop thanos-querier prometheus-global-1 prometheus-global-2 prometheus-local-1 prometheus-local-2 grafana redis jaeger
docker rm thanos-querier prometheus-global-1 prometheus-global-2 prometheus-local-1 prometheus-local-2 grafana redis jaeger
```

#### Start dependencies
- redis
- jaeger
- local-prometheus-1
- local-prometheus-2
- global-prometheus-1
- global-prometheus-2 (HA)
- grafana
```bash
$> docker run -d --name redis -p 6379:6379 redis:latest
$> docker run -d --name jaeger -p 4317:4317 -p 16686:16686 jaegertracing/all-in-one:latest
$> docker run -d --name prometheus-local-1 -p 9091:9090 -v $(pwd)/prometheus-local-1.yml:/etc/prometheus/prometheus.yml prom/prometheus --config.file=/etc/prometheus/prometheus.yml
$> docker run -d --name prometheus-local-2 -p 9092:9090 -v $(pwd)/prometheus-local-2.yml:/etc/prometheus/prometheus.yml prom/prometheus --config.file=/etc/prometheus/prometheus.yml
$> docker run -d --name prometheus-global-1 -p 9090:9090 -v $(pwd)/prometheus-global.yml:/etc/prometheus/prometheus.yml prom/prometheus --config.file=/etc/prometheus/prometheus.yml
$> docker run -d --name prometheus-global-2 -p 9094:9090 -v $(pwd)/prometheus-global.yml:/etc/prometheus/prometheus.yml prom/prometheus --config.file=/etc/prometheus/prometheus.yml
$> docker run -d --name grafana -p 3000:3000 grafana/grafana
```

#### Testing prometheus/ grafana

Terminal 1

```bash
$> go run cmd/main.go -port=8080 -instance=gateway-1
```

Terminal 2

```bash
$> go run cmd/main.go -port=8081 -instance=gateway-2
```
#### Creating grafana dashboards

Open http://localhost:3000, login (admin/admin or your updated password).

“Configuration” > “Data Sources” > Add or edit “Prometheus”:
URL: http://host.docker.internal:9090 (points to prometheus-global-1).

“Save & Test.”

“Dashboards” > “Cache Monitoring”:
Query: rate(cache_hits_total{instance=~".*"}[5m]) and rate(cache_misses_total{instance=~".*"}[5m]).

Datasource: “Prometheus.”

Save if needed.

#### Thanos querier

For true HA, use thanos querier to collect metrics from all prometheus instances, point grafana to thanos endpoint

### Rate limiting

We are using [redis_rate](https://github.com/go-redis/redis_rate) to do per endpoint rate limiting. We can limit at the nginx level as well

To test rate limiting

```bash
$> go run cmd/main.go -port=8080 -instance=gateway-1
```

and in another terminal

```bash
$> bash -c 'for i in {1..15}; do curl -s -o /dev/null -w "%{http_code}\n" "http://localhost:8080/bookings?user_id=123"; done'
200
200
200
200
200
200
200
200
200
200
429
429
429
429
429
```

### Split into services (monorepo)

#### Events

```bash
$> go run cmd/events/main.go -port=8081 -instance=events-1
```

#### Bookings

```bash
$> go run cmd/bookings/main.go -port=8082 -instance=bookings-1
```

#### API Gateway (and Payments)

```bash
$> go run cmd/gateway/main.go -port=8080 -instance=gateway-1
```

## Features

- routing
- zap logging
- jaeger tracing
- open telemetry
- redis API caching
- prometheus (metrics)
- grafana dashboards
- rate limiting (redis_rate)
- split into microservices
