# AI Agentic Crypto Browser - Security Guide

## 🔒 Security Overview

This guide provides comprehensive security measures and best practices for the AI Agentic Crypto Browser, ensuring robust protection against various security threats while maintaining high performance and usability.

## 🛡️ Security Architecture

### Defense in Depth Strategy

```
┌─────────────────────────────────────────────────────────────┐
│                    External Threats                         │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│              Network Security Layer                         │
│  • WAF • DDoS Protection • Rate Limiting • SSL/TLS        │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│            Application Security Layer                       │
│  • Input Validation • Authentication • Authorization       │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│              Data Security Layer                            │
│  • Encryption • Access Controls • Audit Logging           │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│           Infrastructure Security Layer                     │
│  • Container Security • Network Isolation • Monitoring    │
└─────────────────────────────────────────────────────────────┘
```

## 🔐 Authentication and Authorization

### JWT-Based Authentication

```go
type AuthService struct {
    jwtSecret     []byte
    tokenExpiry   time.Duration
    refreshExpiry time.Duration
    userStore     UserStore
    logger        *observability.Logger
}

func (a *AuthService) GenerateTokens(userID uuid.UUID) (*TokenPair, error) {
    // Access token (short-lived)
    accessClaims := &jwt.StandardClaims{
        Subject:   userID.String(),
        ExpiresAt: time.Now().Add(a.tokenExpiry).Unix(),
        IssuedAt:  time.Now().Unix(),
        Issuer:    "ai-agentic-browser",
    }
    
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessTokenString, err := accessToken.SignedString(a.jwtSecret)
    if err != nil {
        return nil, fmt.Errorf("failed to sign access token: %w", err)
    }
    
    // Refresh token (long-lived)
    refreshClaims := &jwt.StandardClaims{
        Subject:   userID.String(),
        ExpiresAt: time.Now().Add(a.refreshExpiry).Unix(),
        IssuedAt:  time.Now().Unix(),
        Issuer:    "ai-agentic-browser",
    }
    
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshTokenString, err := refreshToken.SignedString(a.jwtSecret)
    if err != nil {
        return nil, fmt.Errorf("failed to sign refresh token: %w", err)
    }
    
    return &TokenPair{
        AccessToken:  accessTokenString,
        RefreshToken: refreshTokenString,
        ExpiresIn:    int64(a.tokenExpiry.Seconds()),
    }, nil
}
```

### Role-Based Access Control (RBAC)

```go
type Permission string

const (
    PermissionReadMarketData    Permission = "market:read"
    PermissionWriteMarketData   Permission = "market:write"
    PermissionManageStrategies  Permission = "strategies:manage"
    PermissionViewAnalytics     Permission = "analytics:view"
    PermissionAdminAccess       Permission = "admin:access"
)

type Role struct {
    ID          uuid.UUID    `json:"id"`
    Name        string       `json:"name"`
    Permissions []Permission `json:"permissions"`
}

type User struct {
    ID       uuid.UUID `json:"id"`
    Username string    `json:"username"`
    Email    string    `json:"email"`
    Roles    []Role    `json:"roles"`
    IsActive bool      `json:"is_active"`
}

func (u *User) HasPermission(permission Permission) bool {
    for _, role := range u.Roles {
        for _, perm := range role.Permissions {
            if perm == permission {
                return true
            }
        }
    }
    return false
}

// Authorization middleware
func RequirePermission(permission Permission) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user := GetUserFromContext(r.Context())
            if user == nil {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            
            if !user.HasPermission(permission) {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

## 🔍 Input Validation and Sanitization

### Comprehensive Input Validation

```go
type Validator struct {
    validate *validator.Validate
}

func NewValidator() *Validator {
    v := validator.New()
    
    // Custom validation rules
    v.RegisterValidation("asset_symbol", validateAssetSymbol)
    v.RegisterValidation("strategy_type", validateStrategyType)
    v.RegisterValidation("confidence_range", validateConfidenceRange)
    
    return &Validator{validate: v}
}

func validateAssetSymbol(fl validator.FieldLevel) bool {
    symbol := fl.Field().String()
    // Only allow alphanumeric characters, 2-10 characters long
    matched, _ := regexp.MatchString(`^[A-Z0-9]{2,10}$`, symbol)
    return matched
}

func validateStrategyType(fl validator.FieldLevel) bool {
    strategyType := fl.Field().String()
    validTypes := []string{"trend_following", "mean_reversion", "momentum", "arbitrage"}
    
    for _, valid := range validTypes {
        if strategyType == valid {
            return true
        }
    }
    return false
}

func validateConfidenceRange(fl validator.FieldLevel) bool {
    confidence := fl.Field().Float()
    return confidence >= 0.0 && confidence <= 1.0
}

// Request validation middleware
func ValidationMiddleware(v *Validator) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Validate request size
            if r.ContentLength > MaxRequestSize {
                http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
                return
            }
            
            // Validate content type for POST/PUT requests
            if r.Method == "POST" || r.Method == "PUT" {
                contentType := r.Header.Get("Content-Type")
                if !strings.Contains(contentType, "application/json") {
                    http.Error(w, "Invalid content type", http.StatusUnsupportedMediaType)
                    return
                }
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

### SQL Injection Prevention

```go
// Always use parameterized queries
func (s *StrategyStore) GetStrategiesByType(ctx context.Context, strategyType string) ([]*AdaptiveStrategy, error) {
    query := `
        SELECT id, name, strategy_type, current_parameters, is_active
        FROM market_data.adaptive_strategies 
        WHERE strategy_type = $1 AND is_active = true
        ORDER BY last_adaptation DESC
    `
    
    rows, err := s.db.QueryContext(ctx, query, strategyType)
    if err != nil {
        return nil, fmt.Errorf("failed to query strategies: %w", err)
    }
    defer rows.Close()
    
    var strategies []*AdaptiveStrategy
    for rows.Next() {
        strategy := &AdaptiveStrategy{}
        err := rows.Scan(&strategy.ID, &strategy.Name, &strategy.Type, &strategy.CurrentParameters, &strategy.IsActive)
        if err != nil {
            return nil, fmt.Errorf("failed to scan strategy: %w", err)
        }
        strategies = append(strategies, strategy)
    }
    
    return strategies, nil
}
```

## 🔐 Data Encryption

### Encryption at Rest

```go
type EncryptionService struct {
    key    []byte
    cipher cipher.AEAD
}

func NewEncryptionService(key []byte) (*EncryptionService, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %w", err)
    }
    
    aead, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create AEAD: %w", err)
    }
    
    return &EncryptionService{
        key:    key,
        cipher: aead,
    }, nil
}

func (e *EncryptionService) Encrypt(plaintext []byte) ([]byte, error) {
    nonce := make([]byte, e.cipher.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, fmt.Errorf("failed to generate nonce: %w", err)
    }
    
    ciphertext := e.cipher.Seal(nonce, nonce, plaintext, nil)
    return ciphertext, nil
}

func (e *EncryptionService) Decrypt(ciphertext []byte) ([]byte, error) {
    if len(ciphertext) < e.cipher.NonceSize() {
        return nil, fmt.Errorf("ciphertext too short")
    }
    
    nonce, ciphertext := ciphertext[:e.cipher.NonceSize()], ciphertext[e.cipher.NonceSize():]
    plaintext, err := e.cipher.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt: %w", err)
    }
    
    return plaintext, nil
}

// Encrypt sensitive strategy parameters
func (s *StrategyStore) StoreStrategy(ctx context.Context, strategy *AdaptiveStrategy) error {
    // Encrypt sensitive parameters
    paramData, err := json.Marshal(strategy.CurrentParameters)
    if err != nil {
        return fmt.Errorf("failed to marshal parameters: %w", err)
    }
    
    encryptedParams, err := s.encryption.Encrypt(paramData)
    if err != nil {
        return fmt.Errorf("failed to encrypt parameters: %w", err)
    }
    
    query := `
        INSERT INTO market_data.adaptive_strategies 
        (id, name, strategy_type, encrypted_parameters, is_active)
        VALUES ($1, $2, $3, $4, $5)
    `
    
    _, err = s.db.ExecContext(ctx, query, strategy.ID, strategy.Name, strategy.Type, encryptedParams, strategy.IsActive)
    return err
}
```

### Encryption in Transit

```go
// TLS configuration for production
func NewTLSConfig() *tls.Config {
    return &tls.Config{
        MinVersion:               tls.VersionTLS12,
        CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        },
    }
}

// HTTPS server with security headers
func NewSecureServer(handler http.Handler) *http.Server {
    return &http.Server{
        Addr:         ":8443",
        Handler:      SecurityHeadersMiddleware(handler),
        TLSConfig:    NewTLSConfig(),
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
}

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Security headers
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        
        next.ServeHTTP(w, r)
    })
}
```

## 🚫 Rate Limiting and DDoS Protection

### Advanced Rate Limiting

```go
type RateLimiter struct {
    redis    *redis.Client
    rules    map[string]*RateRule
    logger   *observability.Logger
}

type RateRule struct {
    Requests int           `json:"requests"`
    Window   time.Duration `json:"window"`
    Burst    int           `json:"burst"`
}

func NewRateLimiter(redis *redis.Client) *RateLimiter {
    rules := map[string]*RateRule{
        "api":           {Requests: 100, Window: time.Minute, Burst: 20},
        "auth":          {Requests: 10, Window: time.Minute, Burst: 5},
        "market_data":   {Requests: 1000, Window: time.Minute, Burst: 100},
        "pattern_detect": {Requests: 50, Window: time.Minute, Burst: 10},
    }
    
    return &RateLimiter{
        redis:  redis,
        rules:  rules,
        logger: logger,
    }
}

func (rl *RateLimiter) Allow(ctx context.Context, key string, ruleType string) (bool, error) {
    rule, exists := rl.rules[ruleType]
    if !exists {
        return true, nil // No rule defined, allow
    }
    
    // Use sliding window log algorithm
    now := time.Now()
    windowStart := now.Add(-rule.Window)
    
    pipe := rl.redis.Pipeline()
    
    // Remove old entries
    pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixNano()))
    
    // Count current requests
    pipe.ZCard(ctx, key)
    
    // Add current request
    pipe.ZAdd(ctx, key, &redis.Z{
        Score:  float64(now.UnixNano()),
        Member: fmt.Sprintf("%d", now.UnixNano()),
    })
    
    // Set expiration
    pipe.Expire(ctx, key, rule.Window)
    
    results, err := pipe.Exec(ctx)
    if err != nil {
        return false, fmt.Errorf("rate limit check failed: %w", err)
    }
    
    count := results[1].(*redis.IntCmd).Val()
    
    // Check if within limits
    if count >= int64(rule.Requests) {
        rl.logger.Warn(ctx, "Rate limit exceeded", map[string]interface{}{
            "key":       key,
            "rule_type": ruleType,
            "count":     count,
            "limit":     rule.Requests,
        })
        return false, nil
    }
    
    return true, nil
}

// Rate limiting middleware
func RateLimitMiddleware(limiter *RateLimiter, ruleType string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Use IP address as key
            key := fmt.Sprintf("rate_limit:%s:%s", ruleType, GetClientIP(r))
            
            allowed, err := limiter.Allow(r.Context(), key, ruleType)
            if err != nil {
                http.Error(w, "Internal server error", http.StatusInternalServerError)
                return
            }
            
            if !allowed {
                w.Header().Set("Retry-After", "60")
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

## 🔍 Security Monitoring and Logging

### Comprehensive Audit Logging

```go
type AuditLogger struct {
    logger *observability.Logger
    db     *sql.DB
}

type AuditEvent struct {
    ID        uuid.UUID              `json:"id"`
    UserID    uuid.UUID              `json:"user_id"`
    Action    string                 `json:"action"`
    Resource  string                 `json:"resource"`
    Details   map[string]interface{} `json:"details"`
    IPAddress string                 `json:"ip_address"`
    UserAgent string                 `json:"user_agent"`
    Timestamp time.Time              `json:"timestamp"`
    Success   bool                   `json:"success"`
}

func (al *AuditLogger) LogEvent(ctx context.Context, event *AuditEvent) error {
    // Log to structured logger
    al.logger.Info(ctx, "Audit event", map[string]interface{}{
        "audit_id":   event.ID,
        "user_id":    event.UserID,
        "action":     event.Action,
        "resource":   event.Resource,
        "ip_address": event.IPAddress,
        "success":    event.Success,
    })
    
    // Store in database for compliance
    query := `
        INSERT INTO analytics.audit_events 
        (id, user_id, action, resource, details, ip_address, user_agent, timestamp, success)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
    
    detailsJSON, _ := json.Marshal(event.Details)
    _, err := al.db.ExecContext(ctx, query,
        event.ID, event.UserID, event.Action, event.Resource,
        detailsJSON, event.IPAddress, event.UserAgent, event.Timestamp, event.Success)
    
    return err
}

// Audit middleware
func AuditMiddleware(auditor *AuditLogger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Wrap response writer to capture status
            wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
            
            next.ServeHTTP(wrapped, r)
            
            // Log audit event
            user := GetUserFromContext(r.Context())
            userID := uuid.Nil
            if user != nil {
                userID = user.ID
            }
            
            event := &AuditEvent{
                ID:        uuid.New(),
                UserID:    userID,
                Action:    r.Method,
                Resource:  r.URL.Path,
                Details: map[string]interface{}{
                    "status_code": wrapped.statusCode,
                    "duration_ms": time.Since(start).Milliseconds(),
                },
                IPAddress: GetClientIP(r),
                UserAgent: r.UserAgent(),
                Timestamp: time.Now(),
                Success:   wrapped.statusCode < 400,
            }
            
            auditor.LogEvent(r.Context(), event)
        })
    }
}
```

### Security Monitoring

```go
type SecurityMonitor struct {
    alertManager *AlertManager
    metrics      *SecurityMetrics
    logger       *observability.Logger
}

type SecurityMetrics struct {
    FailedLogins     prometheus.Counter
    RateLimitHits    prometheus.Counter
    SuspiciousEvents prometheus.Counter
    AuthTokenErrors  prometheus.Counter
}

func (sm *SecurityMonitor) DetectAnomalies(ctx context.Context, events []*AuditEvent) {
    for _, event := range events {
        // Detect failed login attempts
        if event.Action == "LOGIN" && !event.Success {
            sm.metrics.FailedLogins.Inc()
            sm.checkFailedLoginThreshold(ctx, event)
        }
        
        // Detect unusual access patterns
        if sm.isUnusualAccess(event) {
            sm.metrics.SuspiciousEvents.Inc()
            sm.alertManager.SendAlert(ctx, &SecurityAlert{
                Type:        "unusual_access",
                Severity:    "medium",
                UserID:      event.UserID,
                IPAddress:   event.IPAddress,
                Description: "Unusual access pattern detected",
                Details:     event.Details,
            })
        }
    }
}

func (sm *SecurityMonitor) checkFailedLoginThreshold(ctx context.Context, event *AuditEvent) {
    // Check if too many failed logins from same IP
    key := fmt.Sprintf("failed_logins:%s", event.IPAddress)
    count, err := sm.redis.Incr(ctx, key).Result()
    if err != nil {
        return
    }
    
    if count == 1 {
        sm.redis.Expire(ctx, key, 15*time.Minute)
    }
    
    if count >= 5 {
        sm.alertManager.SendAlert(ctx, &SecurityAlert{
            Type:        "brute_force",
            Severity:    "high",
            IPAddress:   event.IPAddress,
            Description: fmt.Sprintf("Multiple failed login attempts: %d", count),
        })
    }
}
```

## 🔧 Security Configuration

### Environment-Specific Security

```yaml
# config/security.yaml
security:
  development:
    jwt_secret: "dev-secret-key"
    token_expiry: "1h"
    rate_limiting:
      enabled: false
    encryption:
      enabled: false
    
  production:
    jwt_secret: "${JWT_SECRET}"
    token_expiry: "15m"
    refresh_expiry: "7d"
    rate_limiting:
      enabled: true
      strict_mode: true
    encryption:
      enabled: true
      key_rotation: "30d"
    monitoring:
      security_alerts: true
      audit_logging: true
    
  compliance:
    data_retention: "7y"
    encryption_algorithm: "AES-256-GCM"
    key_management: "HSM"
    audit_requirements:
      - "SOX"
      - "PCI-DSS"
      - "GDPR"
```

### Security Checklist

#### Pre-Deployment Security Checklist

- [ ] **Authentication**
  - [ ] JWT tokens properly configured
  - [ ] Token expiration times set appropriately
  - [ ] Refresh token rotation implemented
  - [ ] Multi-factor authentication available

- [ ] **Authorization**
  - [ ] RBAC properly implemented
  - [ ] Principle of least privilege enforced
  - [ ] API endpoints properly protected
  - [ ] Admin functions restricted

- [ ] **Input Validation**
  - [ ] All inputs validated and sanitized
  - [ ] SQL injection prevention in place
  - [ ] XSS protection implemented
  - [ ] File upload restrictions configured

- [ ] **Encryption**
  - [ ] Data encrypted at rest
  - [ ] TLS 1.2+ for data in transit
  - [ ] Proper key management
  - [ ] Sensitive data identification

- [ ] **Monitoring**
  - [ ] Security logging enabled
  - [ ] Audit trails configured
  - [ ] Anomaly detection active
  - [ ] Alert mechanisms tested

- [ ] **Infrastructure**
  - [ ] Network segmentation implemented
  - [ ] Firewall rules configured
  - [ ] Container security scanned
  - [ ] Dependency vulnerabilities checked

This security guide ensures the AI Agentic Crypto Browser maintains the highest security standards while providing robust protection against various threat vectors.
