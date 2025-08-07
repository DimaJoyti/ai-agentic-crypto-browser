# AI-Agentic Crypto Browser - Production Deployment Guide

## 🚀 Overview

This guide provides comprehensive instructions for deploying the AI-Agentic Crypto Browser in production environments with all enhanced features including advanced trading algorithms, real-time analytics, zero-trust security, and institutional-grade performance optimizations.

## 📋 Prerequisites

### **System Requirements**
- **CPU**: 8+ cores (16+ recommended for high-frequency trading)
- **RAM**: 32GB minimum (64GB+ recommended)
- **Storage**: 1TB SSD (NVMe recommended for low latency)
- **Network**: 10Gbps+ with low latency (<1ms to exchanges)
- **OS**: Ubuntu 22.04 LTS or RHEL 8+

### **Dependencies**
- **Go**: 1.21+ (latest stable)
- **PostgreSQL**: 15+ with TimescaleDB extension
- **Redis**: 7.0+ with Redis Stack modules
- **Docker**: 24.0+ with Docker Compose
- **Kubernetes**: 1.28+ (for container orchestration)
- **Nginx**: 1.24+ (reverse proxy and load balancing)

### **External Services**
- **OpenAI API**: GPT-4 access for AI features
- **Anthropic API**: Claude access for enhanced AI
- **Exchange APIs**: Binance, Coinbase Pro, Kraken, etc.
- **Blockchain RPCs**: Ethereum, Polygon, BSC, Avalanche
- **Monitoring**: Prometheus, Grafana, Jaeger
- **Alerting**: Slack, PagerDuty, email services

## 🏗️ Infrastructure Setup

### **1. Database Configuration**

#### **PostgreSQL with TimescaleDB**
```sql
-- Create database and user
CREATE DATABASE ai_crypto_browser;
CREATE USER ai_browser WITH ENCRYPTED PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE ai_crypto_browser TO ai_browser;

-- Enable TimescaleDB extension
\c ai_crypto_browser;
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Create hypertables for time-series data
SELECT create_hypertable('market_data', 'timestamp');
SELECT create_hypertable('trading_metrics', 'timestamp');
SELECT create_hypertable('performance_metrics', 'timestamp');
```

#### **Redis Configuration**
```conf
# /etc/redis/redis.conf
bind 0.0.0.0
port 6379
requirepass secure_redis_password
maxmemory 8gb
maxmemory-policy allkeys-lru
save 900 1 300 10 60 10000
appendonly yes
appendfsync everysec
tcp-keepalive 60
timeout 300
```

### **2. Environment Configuration**

#### **Production Environment Variables**
```bash
# Application Configuration
export APP_ENV=production
export APP_PORT=8080
export APP_HOST=0.0.0.0

# Database Configuration
export DATABASE_URL="postgres://ai_browser:secure_password@localhost:5432/ai_crypto_browser?sslmode=require"
export DATABASE_READ_REPLICA_URL="postgres://ai_browser:secure_password@replica:5432/ai_crypto_browser?sslmode=require"
export DB_MAX_OPEN_CONNS=100
export DB_MAX_IDLE_CONNS=50
export DB_ENABLE_QUERY_CACHE=true
export DB_CACHE_SIZE=5000

# Redis Configuration
export REDIS_URL="redis://:secure_redis_password@localhost:6379/0"
export REDIS_POOL_SIZE=50
export REDIS_MIN_IDLE_CONNS=20
export REDIS_MAX_IDLE_CONNS=30
export REDIS_ENABLE_METRICS=true
export REDIS_MAX_MEMORY="1gb"

# AI Configuration
export OPENAI_API_KEY="your_openai_api_key"
export ANTHROPIC_API_KEY="your_anthropic_api_key"
export OLLAMA_HOST="http://localhost:11434"

# Web3 Configuration
export ETHEREUM_RPC_URL="https://mainnet.infura.io/v3/your_project_id"
export POLYGON_RPC_URL="https://polygon-mainnet.infura.io/v3/your_project_id"
export BSC_RPC_URL="https://bsc-dataseed.binance.org/"
export AVALANCHE_RPC_URL="https://api.avax.network/ext/bc/C/rpc"

# Security Configuration
export JWT_SECRET="your_jwt_secret_key_256_bits_minimum"
export ENCRYPTION_KEY="your_encryption_key_32_bytes"

# Observability Configuration
export JAEGER_ENDPOINT="http://localhost:14268/api/traces"
export PROMETHEUS_ENDPOINT="http://localhost:9090"

# Trading Configuration
export BINANCE_API_KEY="your_binance_api_key"
export BINANCE_SECRET_KEY="your_binance_secret_key"
export COINBASE_API_KEY="your_coinbase_api_key"
export COINBASE_SECRET_KEY="your_coinbase_secret_key"
```

### **3. Docker Deployment**

#### **Docker Compose Configuration**
```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  # Main Application Services
  api-gateway:
    build:
      context: .
      dockerfile: cmd/api-gateway/Dockerfile
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
    depends_on:
      - postgres
      - redis
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '4.0'
          memory: 8G
        reservations:
          cpus: '2.0'
          memory: 4G

  ai-agent:
    build:
      context: .
      dockerfile: cmd/ai-agent/Dockerfile
    ports:
      - "8082:8082"
    environment:
      - APP_ENV=production
    depends_on:
      - postgres
      - redis
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '8.0'
          memory: 16G
        reservations:
          cpus: '4.0'
          memory: 8G

  browser-service:
    build:
      context: .
      dockerfile: cmd/browser-service/Dockerfile
    ports:
      - "8081:8081"
    environment:
      - APP_ENV=production
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '4.0'
          memory: 8G

  web3-service:
    build:
      context: .
      dockerfile: cmd/web3-service/Dockerfile
    ports:
      - "8083:8083"
    environment:
      - APP_ENV=production
    restart: unless-stopped

  # Database Services
  postgres:
    image: timescale/timescaledb:latest-pg15
    environment:
      POSTGRES_DB: ai_crypto_browser
      POSTGRES_USER: ai_browser
      POSTGRES_PASSWORD: secure_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '4.0'
          memory: 16G

  redis:
    image: redis/redis-stack:latest
    command: redis-server --requirepass secure_redis_password --maxmemory 8gb --maxmemory-policy allkeys-lru
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 8G

  # Monitoring Services
  prometheus:
    image: prom/prometheus:latest
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    ports:
      - "9090:9090"
    restart: unless-stopped

  grafana:
    image: grafana/grafana:latest
    environment:
      GF_SECURITY_ADMIN_PASSWORD: admin_password
    volumes:
      - grafana_data:/var/lib/grafana
      - ./monitoring/grafana/dashboards:/etc/grafana/provisioning/dashboards
    ports:
      - "3000:3000"
    restart: unless-stopped

  jaeger:
    image: jaegertracing/all-in-one:latest
    environment:
      COLLECTOR_OTLP_ENABLED: true
    ports:
      - "16686:16686"
      - "14268:14268"
    restart: unless-stopped

  # Load Balancer
  nginx:
    image: nginx:alpine
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    ports:
      - "80:80"
      - "443:443"
    depends_on:
      - api-gateway
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
  prometheus_data:
  grafana_data:
```

### **4. Kubernetes Deployment**

#### **Kubernetes Manifests**
```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: ai-crypto-browser

---
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: ai-crypto-browser
data:
  APP_ENV: "production"
  DB_MAX_OPEN_CONNS: "100"
  DB_MAX_IDLE_CONNS: "50"
  REDIS_POOL_SIZE: "50"

---
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
  namespace: ai-crypto-browser
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-gateway
  template:
    metadata:
      labels:
        app: api-gateway
    spec:
      containers:
      - name: api-gateway
        image: ai-crypto-browser/api-gateway:latest
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: app-config
        - secretRef:
            name: app-secrets
        resources:
          requests:
            cpu: 2000m
            memory: 4Gi
          limits:
            cpu: 4000m
            memory: 8Gi
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5

---
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: api-gateway-service
  namespace: ai-crypto-browser
spec:
  selector:
    app: api-gateway
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer

---
# k8s/hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: api-gateway-hpa
  namespace: ai-crypto-browser
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: api-gateway
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

## 🔧 Configuration Management

### **1. Production Configuration**
```yaml
# configs/production.yaml (Enhanced)
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 60s
  write_timeout: 60s
  idle_timeout: 300s
  max_header_bytes: 2097152

database:
  max_connections: 100
  max_idle_connections: 50
  max_lifetime: "300s"
  max_idle_time: "300s"
  query_timeout: "30s"
  enable_query_cache: true
  cache_size: 5000
  cache_ttl: "300s"
  enable_read_replica: true
  health_check_interval: "30s"

redis:
  pool_size: 50
  min_idle_connections: 20
  max_idle_connections: 30
  pool_timeout: "10s"
  idle_timeout: "300s"
  max_retries: 5
  enable_metrics: true
  max_memory: "1gb"
  eviction_policy: "allkeys-lru"

security:
  enable_zero_trust: true
  enable_threat_detection: true
  enable_mev_protection: true
  risk_threshold: 0.7
  session_timeout: "30m"
  device_trust_duration: "168h"

trading:
  enable_advanced_orders: true
  enable_cross_chain_arb: true
  enable_liquidity_provision: true
  max_slippage_bps: 50
  max_latency_ms: 100
  risk_tolerance_level: "medium"

analytics:
  enable_real_time_dashboard: true
  enable_predictions: true
  enable_anomaly_detection: true
  update_interval: "1s"
  max_clients: 100
  metrics_retention: "24h"
```

### **2. Monitoring Configuration**
```yaml
# monitoring/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "alert_rules.yml"

scrape_configs:
  - job_name: 'ai-crypto-browser'
    static_configs:
      - targets: ['api-gateway:8080', 'ai-agent:8082', 'browser-service:8081', 'web3-service:8083']
    metrics_path: /metrics
    scrape_interval: 5s

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres:5432']

  - job_name: 'redis'
    static_configs:
      - targets: ['redis:6379']

alerting:
  alertmanagers:
    - static_configs:
        - targets: ['alertmanager:9093']
```

## 🚀 Deployment Steps

### **1. Pre-Deployment Checklist**
```bash
# Verify system requirements
./scripts/check-requirements.sh

# Run security audit
./scripts/security-audit.sh

# Validate configuration
./scripts/validate-config.sh

# Run integration tests
make test-integration

# Build production images
make build-prod

# Run performance benchmarks
./scripts/benchmark.sh
```

### **2. Database Migration**
```bash
# Run database migrations
./scripts/migrate.sh up

# Seed initial data
./scripts/seed-data.sh

# Create indexes for performance
./scripts/create-indexes.sh

# Verify database health
./scripts/db-health-check.sh
```

### **3. Application Deployment**
```bash
# Deploy with Docker Compose
docker-compose -f docker-compose.prod.yml up -d

# Or deploy to Kubernetes
kubectl apply -f k8s/

# Verify deployment
kubectl get pods -n ai-crypto-browser
kubectl get services -n ai-crypto-browser

# Check application health
curl http://localhost:8080/health
curl http://localhost:8082/health
```

### **4. Post-Deployment Verification**
```bash
# Run health checks
./scripts/health-check.sh

# Verify trading algorithms
./scripts/test-trading.sh

# Check security features
./scripts/test-security.sh

# Validate analytics dashboard
./scripts/test-analytics.sh

# Performance validation
./scripts/load-test.sh
```

## 📊 Monitoring & Alerting

### **1. Key Metrics to Monitor**
- **Application**: Response time, throughput, error rate
- **Trading**: Execution latency, slippage, P&L
- **Security**: Threat detection rate, blocked requests
- **Infrastructure**: CPU, memory, disk, network
- **Database**: Query performance, connection pool
- **Cache**: Hit rate, memory usage, evictions

### **2. Alert Rules**
```yaml
# monitoring/alert_rules.yml
groups:
  - name: application
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"

      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"

  - name: trading
    rules:
      - alert: TradingLatencyHigh
        expr: trading_execution_latency_ms > 100
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Trading execution latency too high"

  - name: security
    rules:
      - alert: SecurityThreatDetected
        expr: security_threats_detected_total > 0
        for: 0s
        labels:
          severity: critical
        annotations:
          summary: "Security threat detected"
```

## 🔒 Security Hardening

### **1. Network Security**
```bash
# Configure firewall
ufw enable
ufw allow 22/tcp    # SSH
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw deny 5432/tcp   # PostgreSQL (internal only)
ufw deny 6379/tcp   # Redis (internal only)

# Configure fail2ban
systemctl enable fail2ban
systemctl start fail2ban
```

### **2. SSL/TLS Configuration**
```nginx
# nginx/nginx.conf
server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
    ssl_prefer_server_ciphers off;
    
    location / {
        proxy_pass http://api-gateway:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 🎯 Performance Optimization

### **1. System Tuning**
```bash
# Kernel parameters for high performance
echo 'net.core.rmem_max = 134217728' >> /etc/sysctl.conf
echo 'net.core.wmem_max = 134217728' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_rmem = 4096 87380 134217728' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_wmem = 4096 65536 134217728' >> /etc/sysctl.conf
echo 'vm.swappiness = 1' >> /etc/sysctl.conf
sysctl -p

# CPU governor for performance
echo performance | tee /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor
```

### **2. Application Tuning**
```bash
# Go runtime optimization
export GOGC=100
export GOMAXPROCS=16
export GOMEMLIMIT=32GiB

# Database connection tuning
export DB_MAX_OPEN_CONNS=100
export DB_MAX_IDLE_CONNS=50
export DB_CONN_MAX_LIFETIME=300s

# Redis optimization
export REDIS_POOL_SIZE=50
export REDIS_MAX_MEMORY=8gb
```

This deployment guide provides comprehensive instructions for production deployment with all the enhanced features we've implemented. The configuration ensures optimal performance, security, and reliability for institutional-grade cryptocurrency trading operations.
