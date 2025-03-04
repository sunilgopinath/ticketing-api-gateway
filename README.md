# ticketmaster

I'm going to use hello interview as a basis and implement all the things they talked about 

https://www.hellointerview.com/learn/system-design/problem-breakdowns/ticketmaster

## Getting started

```bash
$> git clone git@github.com:sunilgopinath/ticketing-api-gateway.git
$> cd ticketing-api-gateway
$> go run cmd/main.go
{"level":"info","timestamp":"2025-03-04T09:44:23.733-0800","caller":"cmd/main.go:14","msg":"API Gateway is starting...","port":"8080"}
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


## Features

- routing
- zap logging
- jaeger tracing
- open telemetry
- redis API caching
  