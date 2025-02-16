# Logging Challenge

This project demonstrates a comprehensive logging pipeline using modern observability tools. It features a Go HTTP service with structured logging, which along with Nginx access logs, are collected and processed through Fluent Bit, stored in Loki, and visualized using Grafana dashboards.

## Architecture Overview

```
Nginx/App -> Fluent Bit -> Loki -> Grafana
```

### Components

1. **Source Services**
   - Go HTTP Service (port 8080)
     - Structured logging using zerolog
     - Request ID tracking
     - HTTP method, path, and status code logging
   - Nginx (port 80)
     - Access logs forwarded to Fluent Bit

2. **Log Collection (Fluent Bit)**
   - Collects logs from multiple sources:
     - Tails Go service logs from `/app/logs/app.log`
     - Receives Nginx logs through Forward protocol (port 24224)
   - Parses and processes logs
   - Forwards to Loki with appropriate labels

3. **Log Storage (Loki)**
   - Stores logs with labels for efficient querying
   - Handles log aggregation and indexing
   - Exposes API for Grafana integration

4. **Visualization (Grafana)**
   - Real-time log visualization
   - Pre-configured dashboards
   - Query and analyze logs using LogQL

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Running the Project

1. Clone the repository

2. Start the services:
   ```bash
   docker-compose up -d
   ```

3. Access the services:
   - Go HTTP Service: http://localhost:8080
   - Grafana: http://localhost:3000
   - Nginx: http://localhost:80

## API Endpoints

### Greeting Service
```
GET /?name=YourName
```
Returns a greeting message for the provided name.

## Logging Pipeline Details

### 1. Application Logging

The Go service uses structured logging with the following features:
- Request ID tracking for request tracing
- HTTP method and path logging
- Response status code capture
- Error logging with context

Example log output:
```json
{
  "level": "info",
  "request_id": "uuid",
  "method": "GET",
  "path": "/",
  "status": 200,
  "message": "request completed"
}
```

### 2. Fluent Bit Configuration

Fluent Bit is configured to:
- Tail the application log file
- Receive Nginx logs via forward protocol
- Parse and process logs
- Forward to Loki with appropriate labels

Key configurations:
```ini
[INPUT]
    Name  tail
    Path  /app/logs/app.log
    Tag   http-service

[OUTPUT]
    name        loki
    match       http-service
    host        loki
    port        3100
    labels      app=http-service
```

### 3. Loki Integration

Loki receives logs from Fluent Bit and:
- Labels logs based on source (http-service, nginx)
- Enables efficient log querying
- Provides data source for Grafana

### 4. Grafana Visualization

Grafana is pre-configured with:
- Loki data source
- Dashboards for viewing application and Nginx logs
- LogQL queries for log analysis

## Monitoring and Debugging

1. View logs in Grafana:
   - Navigate to http://localhost:3000
   - Use the pre-configured dashboards
   - Write custom LogQL queries

2. Direct Loki queries:
   - Loki API available at http://localhost:3100
