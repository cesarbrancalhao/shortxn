# Short-xn - URL Shortener Service

A high-performance URL shortening service built with Go, featuring analytics tracking and distributed architecture.

![alt text](img/Example.png)

## Table of Contents

1. [Features](#features)
2. [Architecture](#architecture)
3. [Project Structure](#project-structure)
4. [Prerequisites](#prerequisites)
5. [Quick Start](#quick-start)
6. [API Endpoints](#api-endpoints)
7. [Make Commands](#make-commands) 
8. [Configuration](#configuration)
9. [Monitoring](#monitoring)
10. [Technologies](#technologies)
11. [Usage Guide](#usage-guide)

## Features

- URL shortening with custom length IDs
- Real-time analytics tracking
- Redis caching for fast URL lookups
- Rate limiting and input validation
- Prometheus metrics
- Distributed event processing with RabbitMQ
- PostgreSQL for persistent storage


## Architecture

```
┌─────────────┐         ┌──────────┐        ┌─────────┐
│   API       │ ──────► │  Redis   │        │ Metrics │
│  Service    │         └──────────┘        └─────────┘
└─────────────┘              │
       │                     │             ┌──────────┐
       │                     │        ┌─►  │Consumer 1│
       ▼                     ▼        │    └──────────┘
┌─────────────┐         ┌──────────┐  │    ┌──────────┐
│  PostgreSQL │ ◄────── │ RabbitMQ │──┴─►  │Consumer 2│
└─────────────┘         └──────────┘       └──────────┘
```

## Project Structure

```
shortxn/
├── cmd/
│   ├── api/         # Entrypoint
│   └── consumer/    # Analytics
├── docker/
├── internal/
│   ├── config/
│   ├── domain/      # Models
│   ├── infra/
│   │   ├── postgres/
│   │   ├── redis/
│   │   ├── rabbitmq/
│   │   └── metrics/
│   ├── middleware/  # HTTP middleware
│   ├── ports/       # Interfaces
│   └── service/     # Business logic
└── schemas/         # Database schemas
```

## Prerequisites

- Docker and Docker Compose
- Go 1.22 or later
- Make (optional)

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/yourusername/shortxn.git
cd shortxn
```

2. Start the infrastructure:
```bash
make up
```

3. Run the services:
```bash
# Terminal 1: API Server
make run-api
# or
make start

# Terminal 2: Analytics Consumer
make run-consumer
```

## API Endpoints

### Shorten URL
```bash
POST /shorten
Content-Type: application/json

{
    "url": "https://example.com/very/long/url"
}
```

### Access Shortened URL
```bash
GET /{shortId}
```

### Metrics
```bash
GET /metrics
```

## Make Commands

- `make up`: Start all infrastructure containers
- `make down`: Stop all containers
- `make logs`: View container logs
- `make build`: Build the API and consumer binaries
- `make run-api`: Run the API server
- `make run-consumer`: Run the analytics consumer
- `make clean`: Clean up containers and binaries

## Configuration

Environment variables defaults:
- `POSTGRES_URL`: "postgres://user:password@localhost:5432/shortxn?sslmode=disable"
- `REDIS_ADDR`: "localhost:6379"
- `RABBITMQ_URL`: "amqp://guest:guest@localhost:5672/"
- `SERVER_PORT`: ":8080"
- `RATE_LIMIT`: "100"
- `RATE_BURST`: "50"
- `MAX_URL_LENGTH`: "2048"
- `CACHE_EXPIRATION`: "24"

## Monitoring

- RabbitMQ Management UI: http://localhost:15672
  - Username: guest
  - Password: guest
- Prometheus Metrics: http://localhost:8080/metrics

## Technologies

### Backend
- Go
- Echo
- PostgreSQL
- Redis
- RabbitMQ

### Infrastructure & DevOps
- Docker
- Make
- Prometheus
- Logrus

### Dependencies
- github.com/labstack/echo/v4
- github.com/lib/pq
- github.com/redis/go-redis/v9
- github.com/streadway/amqp
- github.com/prometheus/client_golang
- github.com/sirupsen/logrus
- github.com/kelseyhightower/envconfig
- golang.org/x/time/rate

### Tools
- Git
- GNU Make 4.0+
- Docker Engine 20.10+
- Docker Compose v2.0+

## Usage Guide

### Creating and Using Short URLs

1. Create a short URL:
```bash
curl -X POST http://localhost:8080/shorten \
     -H "Content-Type: application/json" \
     -d '{"url": "https://example.com"}'

# Response:
# {"short_url": "abcde12345=="}
```

2. Access your shortened URL:
   - Using curl:
   ```bash
   curl -L http://localhost:8080/abcde12345==
   ```
   - Or simply open in your browser:
   ```
   http://localhost:8080/abcde12345==
   ```

The service will:
- Automatically redirect to your original URL
- Track analytics for each click
- Cache frequently accessed URLs for faster response
- Record click statistics

### Example Usage

```bash
# 1. Create short URL for a long website address
curl -X POST http://localhost:8080/shorten \
     -H "Content-Type: application/json" \
     -d '{"url": "https://very-long-website-address.com/with/many/segments"}'

# 2. Copy the short_url from the response
# Response: {"short_url": "Ab3CdEfG="} 

# 3. Use the shortened URL
open http://localhost:8080/Ab3CdEfG=
# or
curl -L http://localhost:8080/Ab3CdEfG=
```

### Notes
- The short URLs are case-sensitive
- Each URL generates a unique short code
- The same long URL will always generate the same short code
- Links never expire but are cached for 24 hours