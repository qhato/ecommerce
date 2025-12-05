# Production Configuration Guide

This guide covers production deployment and configuration for the eCommerce platform with full observability and security.

---

## 📋 Prerequisites

### Required Services

- **PostgreSQL 12+** - Primary database
- **Redis 6+** - Caching and rate limiting
- **Prometheus** - Metrics collection
- **Grafana** - Metrics visualization
- **Jaeger** - Distributed tracing (optional but recommended)

### Environment Variables

```bash
# Server Configuration
export PORT=8080
export ENVIRONMENT=production
export SERVICE_NAME=ecommerce-api
export SERVICE_VERSION=1.0.0

# Database
export DATABASE_URL="postgres://user:password@host:5432/ecommerce?sslmode=require"
export DATABASE_MAX_OPEN_CONNS=50
export DATABASE_MAX_IDLE_CONNS=25
export DATABASE_CONN_MAX_LIFETIME=3600

# Redis
export REDIS_URL="redis-host:6379"
export REDIS_PASSWORD="your-redis-password"
export REDIS_DB=0
export REDIS_MAX_RETRIES=3

# JWT Authentication
export JWT_ACCESS_SECRET="generate-strong-random-secret-32-chars-min"
export JWT_REFRESH_SECRET="generate-another-strong-random-secret"
export JWT_ACCESS_TTL=900        # 15 minutes in seconds
export JWT_REFRESH_TTL=604800    # 7 days in seconds

# Tracing
export TRACING_ENABLED=true
export TRACING_EXPORTER=jaeger  # or "otlp"
export JAEGER_ENDPOINT="http://jaeger:14268/api/traces"
export OTLP_ENDPOINT="otel-collector:4317"
export TRACING_SAMPLING_RATE=0.1  # 10% sampling in production

# Rate Limiting
export RATE_LIMIT_ENABLED=true
export RATE_LIMIT_PUBLIC=100         # requests per minute
export RATE_LIMIT_AUTHENTICATED=500  # requests per minute
export RATE_LIMIT_WINDOW=60          # seconds

# CORS
export ALLOWED_ORIGINS="https://app.example.com,https://admin.example.com"

# Logging
export LOG_LEVEL=info        # debug, info, warn, error, fatal
export LOG_FORMAT=json       # json or console
export LOG_OUTPUT=stdout     # stdout, stderr, or file path

# TLS (if terminating TLS in application)
export TLS_ENABLED=false
export TLS_CERT_FILE=/path/to/cert.pem
export TLS_KEY_FILE=/path/to/key.pem
```

---

## 🔐 Security Hardening

### 1. Generate Strong Secrets

```bash
# Generate random secrets for JWT
openssl rand -base64 32  # Access secret
openssl rand -base64 32  # Refresh secret
```

### 2. Configure TLS

**Option A: Application-level TLS**

```go
server := &http.Server{
    Addr:    ":8443",
    Handler: router,
    TLSConfig: &tls.Config{
        MinVersion:               tls.VersionTLS13,
        PreferServerCipherSuites: true,
    },
}

log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
```

**Option B: Reverse Proxy (Recommended)**

Use NGINX, Traefik, or cloud load balancer for TLS termination:

```nginx
server {
    listen 443 ssl http2;
    server_name api.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 3. Database Security

```bash
# Use SSL connections
export DATABASE_URL="postgres://user:pass@host:5432/db?sslmode=require"

# Restrict database user permissions
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO ecommerce_app;
REVOKE CREATE ON SCHEMA public FROM ecommerce_app;
```

### 4. Redis Security

```bash
# Enable password authentication
requirepass your-strong-redis-password

# Disable dangerous commands
rename-command FLUSHDB ""
rename-command FLUSHALL ""
rename-command KEYS ""
rename-command CONFIG ""

# Bind to specific IP
bind 127.0.0.1 ::1
```

---

## 📊 Monitoring Configuration

### Prometheus Configuration

**prometheus.yml**

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    cluster: 'production'
    environment: 'prod'

scrape_configs:
  - job_name: 'ecommerce-api'
    static_configs:
      - targets: ['api-1:8080', 'api-2:8080', 'api-3:8080']
    metrics_path: '/metrics'
    scrape_interval: 10s
    scrape_timeout: 5s

  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'node-exporter'
    static_configs:
      - targets: ['node-exporter:9100']

alerting:
  alertmanagers:
    - static_configs:
        - targets: ['alertmanager:9093']

rule_files:
  - '/etc/prometheus/rules/*.yml'
```

### Alerting Rules

**alerts.yml**

```yaml
groups:
  - name: ecommerce_alerts
    interval: 30s
    rules:
      # High error rate
      - alert: HighErrorRate
        expr: rate(ecommerce_http_errors_total[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value | humanizePercentage }} (> 5%)"

      # High latency
      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(ecommerce_http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High API latency"
          description: "95th percentile latency is {{ $value }}s"

      # Low success rate
      - alert: LowSuccessRate
        expr: rate(ecommerce_http_requests_total{status=~"2.."}[5m]) / rate(ecommerce_http_requests_total[5m]) < 0.95
        for: 10m
        labels:
          severity: critical
        annotations:
          summary: "Low success rate"
          description: "Success rate is {{ $value | humanizePercentage }}"

      # Database connection pool exhaustion
      - alert: DatabaseConnectionPoolNearLimit
        expr: ecommerce_database_connections_open / ecommerce_database_connections_max > 0.9
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Database connection pool nearly exhausted"
          description: "{{ $value | humanizePercentage }} of connections in use"

      # High cache miss rate
      - alert: HighCacheMissRate
        expr: rate(ecommerce_cache_misses_total[5m]) / (rate(ecommerce_cache_hits_total[5m]) + rate(ecommerce_cache_misses_total[5m])) > 0.5
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "High cache miss rate"
          description: "Cache miss rate is {{ $value | humanizePercentage }}"
```

### Grafana Dashboards

**Import these dashboard IDs:**

1. **Go Application Metrics (12486)**
   - Memory usage
   - Goroutines
   - GC stats

2. **Prometheus Stats (13408)**
   - Prometheus performance
   - Target status

3. **Custom Business Dashboard**

```json
{
  "dashboard": {
    "title": "ECommerce Business Metrics",
    "panels": [
      {
        "title": "Orders Per Minute",
        "targets": [{
          "expr": "rate(ecommerce_orders_created_total[1m]) * 60"
        }]
      },
      {
        "title": "Revenue Per Hour",
        "targets": [{
          "expr": "rate(ecommerce_order_value_dollars_sum[1h]) * 3600"
        }]
      },
      {
        "title": "Payment Success Rate",
        "targets": [{
          "expr": "rate(ecommerce_payments_processed_total[5m]) / (rate(ecommerce_payments_processed_total[5m]) + rate(ecommerce_payments_failed_total[5m]))"
        }]
      }
    ]
  }
}
```

---

## 🚀 Deployment

### Docker

**Dockerfile (Production)**

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
    -o /app/bin/ecommerce-api cmd/admin/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/ecommerce-api .

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health/live || exit 1

ENTRYPOINT ["./ecommerce-api"]
```

**Build and Run**

```bash
# Build
docker build -t ecommerce-api:1.0.0 .

# Run
docker run -d \
    --name ecommerce-api \
    -p 8080:8080 \
    -e DATABASE_URL="postgres://..." \
    -e REDIS_URL="redis:6379" \
    -e JWT_ACCESS_SECRET="..." \
    -e JWT_REFRESH_SECRET="..." \
    ecommerce-api:1.0.0
```

### Kubernetes

**deployment.yaml**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ecommerce-api
  labels:
    app: ecommerce-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ecommerce-api
  template:
    metadata:
      labels:
        app: ecommerce-api
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
        prometheus.io/path: "/metrics"
    spec:
      containers:
      - name: api
        image: ecommerce-api:1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: ENVIRONMENT
          value: "production"
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: ecommerce-secrets
              key: database-url
        - name: REDIS_URL
          valueFrom:
            configMapKeyRef:
              name: ecommerce-config
              key: redis-url
        - name: JWT_ACCESS_SECRET
          valueFrom:
            secretKeyRef:
              name: ecommerce-secrets
              key: jwt-access-secret
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"

---
apiVersion: v1
kind: Service
metadata:
  name: ecommerce-api
spec:
  selector:
    app: ecommerce-api
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: ecommerce-config
data:
  redis-url: "redis:6379"
  log-level: "info"
  tracing-enabled: "true"

---
apiVersion: v1
kind: Secret
metadata:
  name: ecommerce-secrets
type: Opaque
data:
  database-url: <base64-encoded-url>
  jwt-access-secret: <base64-encoded-secret>
  jwt-refresh-secret: <base64-encoded-secret>
```

**Deploy**

```bash
kubectl apply -f deployment.yaml
kubectl rollout status deployment/ecommerce-api
```

### Horizontal Pod Autoscaler

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: ecommerce-api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: ecommerce-api
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

---

## 🔍 Observability in Production

### Jaeger Configuration

**Recommended setup:** Use Jaeger with Elasticsearch backend

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: jaeger-config
data:
  jaeger.yml: |
    reporter:
      queueSize: 1000
      bufferFlushInterval: 1s
    sampler:
      type: probabilistic
      param: 0.1  # 10% sampling
```

### Log Aggregation

**Use centralized logging:**

- **ELK Stack** (Elasticsearch, Logstash, Kibana)
- **Loki** (Grafana Loki)
- **Cloud Solutions** (CloudWatch, Stackdriver, DataDog)

**Structured logging makes this easy:**

```json
{
  "level": "info",
  "timestamp": "2025-12-01T10:30:00Z",
  "trace_id": "abc123",
  "span_id": "def456",
  "service": "ecommerce-api",
  "method": "POST",
  "path": "/api/v1/orders",
  "status": 201,
  "duration_ms": 45,
  "user_id": 12345,
  "order_id": 67890,
  "msg": "Order created successfully"
}
```

---

## 🎯 Performance Tuning

### Database Connection Pool

```go
db.SetMaxOpenConns(50)              // Maximum open connections
db.SetMaxIdleConns(25)              // Maximum idle connections
db.SetConnMaxLifetime(time.Hour)    // Maximum lifetime of connection
db.SetConnMaxIdleTime(10 * time.Minute) // Maximum idle time
```

### Redis Configuration

```go
redis.NewClient(&redis.Options{
    Addr:         "redis:6379",
    Password:     os.Getenv("REDIS_PASSWORD"),
    DB:           0,
    PoolSize:     50,
    MinIdleConns: 10,
    MaxRetries:   3,
    DialTimeout:  5 * time.Second,
    ReadTimeout:  3 * time.Second,
    WriteTimeout: 3 * time.Second,
})
```

### HTTP Server

```go
server := &http.Server{
    Addr:         ":8080",
    Handler:      router,
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
    IdleTimeout:  60 * time.Second,
    MaxHeaderBytes: 1 << 20, // 1MB
}
```

---

## 📝 Maintenance

### Backup Strategy

```bash
# Database backups
pg_dump -h postgres -U user -d ecommerce > backup-$(date +%Y%m%d).sql

# Automated daily backups
0 2 * * * /usr/local/bin/backup-db.sh
```

### Log Rotation

```yaml
# logrotate config
/var/log/ecommerce/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 ecommerce ecommerce
    sharedscripts
    postrotate
        /usr/bin/killall -SIGUSR1 ecommerce-api
    endscript
}
```

### Database Migrations

```bash
# Run migrations
psql $DATABASE_URL < migrations/20251201_add_feature.sql

# Rollback (if needed)
psql $DATABASE_URL < migrations/20251201_add_feature.down.sql
```

---

## 🚨 Incident Response

### Health Check Endpoints

```bash
# Check overall health
curl https://api.example.com/health

# Check liveness (restart if fails)
curl https://api.example.com/health/live

# Check readiness (remove from load balancer if fails)
curl https://api.example.com/health/ready
```

### Metrics for Debugging

```bash
# Check error rate
curl https://api.example.com/metrics | grep http_errors_total

# Check latency
curl https://api.example.com/metrics | grep http_request_duration_seconds

# Check database connections
curl https://api.example.com/metrics | grep database_connections
```

### Tracing for Root Cause Analysis

1. Find failing requests in Jaeger
2. Analyze trace timeline
3. Identify slow spans
4. Check for errors in span events

---

## ✅ Production Checklist

- [ ] Strong JWT secrets generated
- [ ] TLS/SSL configured
- [ ] Database SSL enabled
- [ ] Redis password set
- [ ] Rate limiting enabled
- [ ] Security headers configured
- [ ] CORS origins whitelisted
- [ ] Prometheus scraping configured
- [ ] Grafana dashboards imported
- [ ] Alerting rules configured
- [ ] Log aggregation setup
- [ ] Backup strategy implemented
- [ ] Monitoring dashboards created
- [ ] Incident response plan documented
- [ ] Load testing completed
- [ ] Health checks verified
- [ ] Horizontal scaling tested

---

## 📚 Additional Resources

- [Prometheus Best Practices](https://prometheus.io/docs/practices/)
- [Grafana Dashboard Best Practices](https://grafana.com/docs/grafana/latest/best-practices/)
- [OpenTelemetry Documentation](https://opentelemetry.io/docs/)
- [JWT Security Best Practices](https://datatracker.ietf.org/doc/html/rfc8725)
- [Go Security Checklist](https://go.dev/doc/security/checklist)
