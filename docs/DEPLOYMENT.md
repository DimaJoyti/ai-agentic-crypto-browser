# AI Agentic Crypto Browser - Deployment Guide

This guide provides comprehensive instructions for deploying the AI Agentic Crypto Browser in various environments.

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)
- PostgreSQL 15+
- Redis 7+
- 4GB+ RAM
- 2+ CPU cores

### Local Development

```bash
# Clone the repository
git clone https://github.com/your-org/ai-agentic-crypto-browser.git
cd ai-agentic-crypto-browser

# Build and run locally
go build ./cmd/ai-agent/
./ai-agent

# Or run with Docker Compose
cd deployments/docker
docker-compose up -d
```

## 🐳 Docker Deployment

### Single Container

```bash
# Build the image
docker build -f deployments/docker/Dockerfile -t ai-agentic-crypto-browser .

# Run the container
docker run -d \
  --name ai-agent \
  -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@host:5432/db" \
  -e REDIS_URL="redis://host:6379" \
  ai-agentic-crypto-browser
```

### Docker Compose (Recommended)

```bash
cd deployments/docker
docker-compose up -d
```

This will start:
- AI Agent service (port 8080)
- PostgreSQL database (port 5432)
- Redis cache (port 6379)
- Prometheus monitoring (port 9090)
- Grafana dashboard (port 3000)
- Nginx reverse proxy (port 80/443)

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENV` | Environment (development/production) | `development` |
| `LOG_LEVEL` | Logging level (debug/info/warn/error) | `info` |
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `REDIS_URL` | Redis connection string | Required |
| `BROWSER_HEADLESS` | Run browser in headless mode | `true` |
| `BROWSER_TIMEOUT` | Browser operation timeout | `30s` |
| `AI_MODEL_TIMEOUT` | AI model processing timeout | `60s` |
| `MARKET_ADAPTATION_ENABLED` | Enable market pattern adaptation | `true` |
| `PATTERN_DETECTION_WINDOW` | Pattern detection time window | `7d` |
| `ADAPTATION_THRESHOLD` | Adaptation confidence threshold | `0.7` |
| `REAL_TIME_ADAPTATION` | Enable real-time adaptation | `true` |

## ☸️ Kubernetes Deployment

### Prerequisites

- Kubernetes cluster (1.20+)
- kubectl configured
- Helm 3+ (optional)
- Ingress controller (nginx recommended)
- Cert-manager (for SSL)

### Deploy to Kubernetes

```bash
# Create namespace
kubectl create namespace ai-browser

# Apply the deployment
kubectl apply -f deployments/k8s/ai-agent-deployment.yaml

# Check deployment status
kubectl get pods -n ai-browser
kubectl get services -n ai-browser
```

### Scaling

```bash
# Manual scaling
kubectl scale deployment ai-agentic-crypto-browser --replicas=5 -n ai-browser

# Auto-scaling is configured via HPA (3-10 replicas)
kubectl get hpa -n ai-browser
```

### Monitoring

```bash
# Check logs
kubectl logs -f deployment/ai-agentic-crypto-browser -n ai-browser

# Check metrics
kubectl port-forward service/ai-agent-service 8080:80 -n ai-browser
curl http://localhost:8080/metrics
```

## 🌐 Production Deployment

### Infrastructure Requirements

#### Minimum Requirements
- **CPU**: 4 cores
- **Memory**: 8GB RAM
- **Storage**: 50GB SSD
- **Network**: 1Gbps

#### Recommended for High Load
- **CPU**: 8+ cores
- **Memory**: 16GB+ RAM
- **Storage**: 100GB+ NVMe SSD
- **Network**: 10Gbps
- **Load Balancer**: HAProxy/nginx
- **Database**: PostgreSQL cluster
- **Cache**: Redis cluster

### Security Configuration

#### SSL/TLS Setup

```bash
# Generate SSL certificates (Let's Encrypt)
certbot certonly --webroot -w /var/www/html -d api.your-domain.com

# Update nginx configuration
cp deployments/docker/nginx.conf /etc/nginx/sites-available/ai-agent
# Uncomment SSL configuration section
```

#### Firewall Rules

```bash
# Allow only necessary ports
ufw allow 22/tcp    # SSH
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw deny 5432/tcp   # PostgreSQL (internal only)
ufw deny 6379/tcp   # Redis (internal only)
ufw enable
```

#### Database Security

```sql
-- Create dedicated database user
CREATE USER ai_agent WITH PASSWORD 'secure_password';
GRANT CONNECT ON DATABASE ai_browser TO ai_agent;
GRANT USAGE ON SCHEMA ai_data, market_data, user_data, analytics TO ai_agent;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA ai_data, market_data, user_data, analytics TO ai_agent;

-- Enable SSL
ALTER SYSTEM SET ssl = on;
ALTER SYSTEM SET ssl_cert_file = '/path/to/server.crt';
ALTER SYSTEM SET ssl_key_file = '/path/to/server.key';
```

### Performance Optimization

#### Database Optimization

```sql
-- Optimize PostgreSQL settings
ALTER SYSTEM SET shared_buffers = '2GB';
ALTER SYSTEM SET effective_cache_size = '6GB';
ALTER SYSTEM SET maintenance_work_mem = '512MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET default_statistics_target = 100;
ALTER SYSTEM SET random_page_cost = 1.1;
ALTER SYSTEM SET effective_io_concurrency = 200;

-- Reload configuration
SELECT pg_reload_conf();
```

#### Redis Optimization

```bash
# Redis configuration
echo "maxmemory 4gb" >> /etc/redis/redis.conf
echo "maxmemory-policy allkeys-lru" >> /etc/redis/redis.conf
echo "save 900 1" >> /etc/redis/redis.conf
echo "save 300 10" >> /etc/redis/redis.conf
echo "save 60 10000" >> /etc/redis/redis.conf
```

#### Application Optimization

```yaml
# config.yaml
server:
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 120s
  max_header_bytes: 1048576

ai:
  max_concurrent_requests: 100
  request_queue_size: 1000
  worker_pool_size: 50

market_adaptation:
  pattern_cache_size: 10000
  strategy_cache_ttl: 3600s
  performance_cache_ttl: 1800s
```

## 📊 Monitoring and Observability

### Metrics Collection

The application exposes Prometheus metrics at `/metrics`:

- **Request metrics**: Duration, count, errors
- **AI metrics**: Processing time, model performance
- **Market metrics**: Pattern detection, strategy adaptation
- **System metrics**: Memory, CPU, goroutines

### Grafana Dashboards

Access Grafana at `http://localhost:3000` (admin/admin):

1. **AI Agent Overview**: System health, request rates, errors
2. **Market Analysis**: Pattern detection, strategy performance
3. **User Behavior**: Interaction patterns, learning metrics
4. **Infrastructure**: Database, Redis, system resources

### Alerting Rules

```yaml
# prometheus-alerts.yml
groups:
- name: ai-agent-alerts
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: High error rate detected

  - alert: HighMemoryUsage
    expr: process_resident_memory_bytes / 1024 / 1024 > 1500
    for: 10m
    labels:
      severity: warning
    annotations:
      summary: High memory usage

  - alert: DatabaseConnectionFailure
    expr: up{job="postgres"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: Database connection failed
```

### Log Management

```bash
# Centralized logging with ELK stack
docker run -d \
  --name elasticsearch \
  -p 9200:9200 \
  -e "discovery.type=single-node" \
  elasticsearch:7.17.0

docker run -d \
  --name logstash \
  -p 5044:5044 \
  logstash:7.17.0

docker run -d \
  --name kibana \
  -p 5601:5601 \
  -e "ELASTICSEARCH_HOSTS=http://elasticsearch:9200" \
  kibana:7.17.0
```

## 🔧 Maintenance

### Database Maintenance

```sql
-- Regular maintenance tasks
VACUUM ANALYZE;
REINDEX DATABASE ai_browser;

-- Clean old data
DELETE FROM user_data.behavior_events WHERE timestamp < NOW() - INTERVAL '90 days';
DELETE FROM analytics.api_usage WHERE timestamp < NOW() - INTERVAL '30 days';
```

### Backup Strategy

```bash
#!/bin/bash
# backup.sh

# Database backup
pg_dump -h localhost -U postgres ai_browser | gzip > backup_$(date +%Y%m%d_%H%M%S).sql.gz

# Redis backup
redis-cli --rdb backup_redis_$(date +%Y%m%d_%H%M%S).rdb

# Application data backup
tar -czf app_data_$(date +%Y%m%d_%H%M%S).tar.gz /app/data/
```

### Health Checks

```bash
#!/bin/bash
# health-check.sh

# Check application health
curl -f http://localhost:8080/health || exit 1

# Check database connectivity
pg_isready -h localhost -p 5432 -U postgres || exit 1

# Check Redis connectivity
redis-cli ping || exit 1

# Check disk space
df -h | awk '$5 > 80 {print "Disk usage high: " $0; exit 1}'
```

## 🚨 Troubleshooting

### Common Issues

#### High Memory Usage
```bash
# Check memory usage
docker stats
kubectl top pods -n ai-browser

# Optimize garbage collection
export GOGC=100
export GOMEMLIMIT=2GiB
```

#### Database Connection Issues
```bash
# Check connection limits
SELECT count(*) FROM pg_stat_activity;
ALTER SYSTEM SET max_connections = 200;

# Check connection pooling
SELECT * FROM pg_stat_database WHERE datname = 'ai_browser';
```

#### Performance Issues
```bash
# Profile the application
go tool pprof http://localhost:8080/debug/pprof/profile

# Check slow queries
SELECT query, mean_time, calls FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;
```

### Support

- **Documentation**: Check the `docs/` directory
- **Issues**: Report bugs on GitHub
- **Monitoring**: Use Grafana dashboards for insights
- **Logs**: Check application logs for detailed error information

## 📋 Deployment Checklist

### Pre-deployment
- [ ] Infrastructure provisioned
- [ ] SSL certificates configured
- [ ] Database initialized
- [ ] Environment variables set
- [ ] Security configurations applied
- [ ] Monitoring setup complete

### Deployment
- [ ] Application deployed
- [ ] Health checks passing
- [ ] Database migrations applied
- [ ] Cache warmed up
- [ ] Load balancer configured
- [ ] DNS records updated

### Post-deployment
- [ ] Monitoring alerts configured
- [ ] Backup strategy implemented
- [ ] Performance baseline established
- [ ] Documentation updated
- [ ] Team notified
- [ ] Rollback plan prepared

## 🔄 Updates and Rollbacks

### Rolling Updates
```bash
# Kubernetes rolling update
kubectl set image deployment/ai-agentic-crypto-browser ai-agent=ai-agentic-crypto-browser:v1.1.0 -n ai-browser

# Docker Compose update
docker-compose pull
docker-compose up -d
```

### Rollback
```bash
# Kubernetes rollback
kubectl rollout undo deployment/ai-agentic-crypto-browser -n ai-browser

# Docker rollback
docker-compose down
docker-compose up -d --force-recreate
```

This deployment guide ensures a robust, scalable, and maintainable production deployment of the AI Agentic Crypto Browser.
