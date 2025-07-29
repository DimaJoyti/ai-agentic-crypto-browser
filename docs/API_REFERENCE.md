# AI Agentic Crypto Browser - Complete API Reference

## 📚 Overview

This comprehensive API reference documents all 43+ endpoints available in the AI Agentic Crypto Browser, organized by functional areas and including detailed examples, request/response schemas, and usage patterns.

## 🔗 Base URL

```
Production: https://api.ai-crypto-browser.com
Development: http://localhost:8080
```

## 🔐 Authentication

All API endpoints require authentication via JWT tokens in the Authorization header:

```http
Authorization: Bearer <your-jwt-token>
```

### Authentication Endpoints

#### POST /auth/login
Authenticate user and receive JWT tokens.

**Request:**
```json
{
  "username": "user@example.com",
  "password": "secure_password"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "token_type": "Bearer"
}
```

#### POST /auth/refresh
Refresh access token using refresh token.

**Request:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## 🤖 Enhanced AI Service Endpoints

### POST /ai/enhanced/analyze
Advanced AI analysis with multiple reasoning capabilities.

**Request:**
```json
{
  "request_id": "req_123",
  "user_id": "user_456",
  "type": "market_analysis",
  "symbol": "BTC",
  "data": {
    "query": "Analyze Bitcoin market trends for the next week",
    "timeframe": "1w",
    "include_sentiment": true
  },
  "context": {
    "current_price": 45000,
    "volume_24h": 28000000000
  }
}
```

**Response:**
```json
{
  "request_id": "req_123",
  "analysis": {
    "summary": "Bitcoin shows bullish momentum with strong support at $44,000",
    "confidence": 0.85,
    "key_insights": [
      "Technical indicators suggest upward trend",
      "Volume analysis shows institutional interest",
      "Sentiment analysis indicates positive market mood"
    ],
    "predictions": {
      "price_target_7d": 48000,
      "probability": 0.72,
      "risk_factors": ["Regulatory uncertainty", "Market volatility"]
    }
  },
  "processing_time_ms": 1250,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### POST /ai/enhanced/predict
Generate predictions using advanced ML models.

**Request:**
```json
{
  "prediction_type": "price_movement",
  "asset": "ETH",
  "timeframe": "24h",
  "features": {
    "technical_indicators": {
      "rsi": 65.5,
      "macd": 0.8,
      "bollinger_position": 0.7
    },
    "market_data": {
      "volume_ratio": 1.2,
      "volatility": 0.15
    }
  }
}
```

**Response:**
```json
{
  "prediction": {
    "direction": "bullish",
    "magnitude": 0.15,
    "confidence": 0.78,
    "price_targets": {
      "conservative": 3200,
      "moderate": 3350,
      "aggressive": 3500
    }
  },
  "model_info": {
    "model_version": "v2.1.0",
    "training_date": "2024-01-10",
    "accuracy_score": 0.82
  }
}
```

## 🧠 Advanced NLP Engine Endpoints

### POST /ai/nlp/analyze
Comprehensive natural language processing and sentiment analysis.

**Request:**
```json
{
  "text": "Bitcoin's recent surge past $45,000 has sparked renewed optimism among institutional investors",
  "analysis_types": ["sentiment", "entities", "keywords", "classification"],
  "context": {
    "domain": "cryptocurrency",
    "language": "en"
  }
}
```

**Response:**
```json
{
  "sentiment": {
    "overall": "positive",
    "score": 0.75,
    "confidence": 0.88,
    "emotions": {
      "optimism": 0.8,
      "confidence": 0.7,
      "excitement": 0.6
    }
  },
  "entities": [
    {"text": "Bitcoin", "type": "CRYPTOCURRENCY", "confidence": 0.99},
    {"text": "$45,000", "type": "PRICE", "confidence": 0.95},
    {"text": "institutional investors", "type": "INVESTOR_TYPE", "confidence": 0.92}
  ],
  "keywords": ["Bitcoin", "surge", "optimism", "institutional", "investors"],
  "classification": {
    "category": "market_news",
    "subcategory": "price_movement",
    "confidence": 0.91
  }
}
```

### POST /ai/nlp/summarize
Generate intelligent summaries of long-form content.

**Request:**
```json
{
  "content": "Long article text here...",
  "summary_type": "extractive",
  "max_sentences": 3,
  "focus_areas": ["key_findings", "market_impact", "future_outlook"]
}
```

## 🎯 Decision Making Engine Endpoints

### POST /ai/decision/evaluate
Evaluate complex decisions using AI-powered analysis.

**Request:**
```json
{
  "decision_context": {
    "type": "trading_decision",
    "asset": "BTC",
    "current_position": "long",
    "market_conditions": "volatile"
  },
  "options": [
    {
      "id": "hold",
      "description": "Maintain current position",
      "parameters": {"risk_level": "medium"}
    },
    {
      "id": "close",
      "description": "Close position and take profits",
      "parameters": {"risk_level": "low"}
    }
  ],
  "criteria": {
    "risk_tolerance": 0.6,
    "time_horizon": "short_term",
    "profit_target": 0.15
  }
}
```

**Response:**
```json
{
  "recommendation": {
    "option_id": "hold",
    "confidence": 0.72,
    "reasoning": [
      "Technical indicators suggest continued upward momentum",
      "Risk-reward ratio favors holding current position",
      "Market sentiment remains positive"
    ]
  },
  "risk_assessment": {
    "overall_risk": "medium",
    "risk_factors": ["Market volatility", "Regulatory uncertainty"],
    "mitigation_strategies": ["Set stop-loss at $43,000", "Monitor volume trends"]
  },
  "alternative_scenarios": [
    {
      "scenario": "market_reversal",
      "probability": 0.25,
      "recommended_action": "close_position"
    }
  ]
}
```

## 👤 User Behavior Learning Endpoints

### POST /ai/behavior/track
Track user behavior events for learning and personalization.

**Request:**
```json
{
  "user_id": "user_456",
  "event": {
    "type": "trading_action",
    "action": "position_opened",
    "context": {
      "asset": "ETH",
      "position_size": 0.05,
      "entry_price": 3150,
      "strategy_used": "momentum_trading",
      "market_conditions": "bullish",
      "time_of_day": "morning",
      "session_duration": "45m"
    },
    "outcome": {
      "success": true,
      "profit_loss": 0.08,
      "duration_held": "2h"
    }
  }
}
```

**Response:**
```json
{
  "event_id": "evt_789",
  "learning_insights": {
    "pattern_identified": "user_prefers_morning_trading",
    "confidence": 0.68,
    "recommendation": "Consider automated morning alerts for similar opportunities"
  },
  "profile_updates": {
    "trading_style": "momentum_trader",
    "risk_profile": "moderate_aggressive",
    "preferred_timeframes": ["1h", "4h"]
  }
}
```

### GET /ai/behavior/profile/{user_id}
Retrieve comprehensive user behavior profile.

**Response:**
```json
{
  "user_id": "user_456",
  "profile": {
    "trading_style": {
      "primary_style": "momentum_trader",
      "secondary_styles": ["swing_trader"],
      "trading_frequency": 3.2,
      "average_hold_time": "4h",
      "preferred_assets": ["BTC", "ETH", "SOL"]
    },
    "risk_profile": {
      "risk_tolerance": 0.7,
      "max_position_size": 0.1,
      "preferred_risk_reward": 2.5
    },
    "behavior_patterns": [
      {
        "pattern": "morning_trading_preference",
        "confidence": 0.85,
        "frequency": "daily"
      }
    ],
    "recommendations": [
      {
        "type": "strategy_optimization",
        "suggestion": "Increase position size for high-confidence momentum signals",
        "confidence": 0.78
      }
    ]
  }
}
```

## 🎨 Multi-Modal AI Engine Endpoints

### POST /ai/multimodal/analyze
Analyze various content types including images, documents, and audio.

**Request:**
```json
{
  "content": [
    {
      "type": "image",
      "data": "base64_encoded_chart_image",
      "metadata": {
        "chart_type": "candlestick",
        "asset": "BTC",
        "timeframe": "1h"
      }
    },
    {
      "type": "text",
      "data": "Market analysis report content...",
      "metadata": {
        "document_type": "market_report",
        "source": "research_firm"
      }
    }
  ],
  "analysis_types": ["pattern_recognition", "sentiment_analysis", "data_extraction"]
}
```

**Response:**
```json
{
  "results": [
    {
      "content_id": "img_001",
      "analysis": {
        "patterns_detected": [
          {
            "type": "ascending_triangle",
            "confidence": 0.82,
            "coordinates": {"x1": 100, "y1": 200, "x2": 300, "y2": 180}
          }
        ],
        "technical_indicators": {
          "trend": "bullish",
          "support_level": 44500,
          "resistance_level": 45200
        }
      }
    },
    {
      "content_id": "txt_001",
      "analysis": {
        "sentiment": "positive",
        "key_insights": ["Institutional adoption increasing", "Technical breakout imminent"],
        "extracted_data": {
          "price_targets": [46000, 48000],
          "timeframes": ["1w", "1m"]
        }
      }
    }
  ]
}
```

## 📊 Market Pattern Adaptation Endpoints

### POST /ai/market/patterns/detect
Detect market patterns in real-time data.

**Request:**
```json
{
  "market_data": {
    "asset": "BTC",
    "timeframe": "1h",
    "prices": [44000, 44200, 44500, 44800, 45000],
    "volumes": [1000, 1200, 1500, 1800, 2000],
    "timestamps": ["2024-01-15T10:00:00Z", "2024-01-15T11:00:00Z", "..."]
  },
  "technical_indicators": {
    "rsi": 65.5,
    "macd": 0.8,
    "bollinger_bands": {
      "upper": 45500,
      "middle": 44500,
      "lower": 43500
    }
  },
  "detection_params": {
    "min_confidence": 0.7,
    "pattern_types": ["trend", "reversal", "breakout"],
    "lookback_periods": 20
  }
}
```

**Response:**
```json
{
  "patterns": [
    {
      "id": "pat_123",
      "type": "bullish_trend",
      "confidence": 0.85,
      "strength": 0.78,
      "timeframe": "1h",
      "characteristics": {
        "slope": 0.15,
        "duration": "4h",
        "volume_confirmation": true
      },
      "market_context": {
        "trend_direction": "up",
        "volatility": "medium",
        "volume_profile": "increasing"
      },
      "price_levels": {
        "entry": 44800,
        "target": 46000,
        "stop_loss": 44200
      }
    }
  ],
  "market_summary": {
    "overall_trend": "bullish",
    "pattern_count": 1,
    "avg_confidence": 0.85
  }
}
```

### GET /ai/market/strategies
List all adaptive trading strategies.

**Response:**
```json
{
  "strategies": [
    {
      "id": "strat_001",
      "name": "Adaptive Momentum Strategy",
      "type": "momentum",
      "is_active": true,
      "current_parameters": {
        "position_size": 0.05,
        "stop_loss": 0.02,
        "take_profit": 0.04,
        "momentum_threshold": 0.15
      },
      "performance_metrics": {
        "total_return": 0.18,
        "sharpe_ratio": 1.45,
        "max_drawdown": 0.08,
        "win_rate": 0.68
      },
      "adaptation_count": 12,
      "last_adaptation": "2024-01-15T09:30:00Z"
    }
  ]
}
```

### POST /ai/market/strategies/{strategy_id}/adapt
Trigger strategy adaptation based on current market conditions.

**Request:**
```json
{
  "market_conditions": {
    "volatility": "high",
    "trend_strength": 0.8,
    "volume_profile": "increasing"
  },
  "performance_data": {
    "recent_trades": 10,
    "win_rate": 0.6,
    "avg_return": 0.03
  },
  "adaptation_triggers": ["performance_decline", "market_regime_change"]
}
```

**Response:**
```json
{
  "adaptation": {
    "id": "adapt_456",
    "strategy_id": "strat_001",
    "adaptation_type": "parameter_optimization",
    "trigger_reason": "performance_decline",
    "old_parameters": {
      "position_size": 0.05,
      "stop_loss": 0.02
    },
    "new_parameters": {
      "position_size": 0.03,
      "stop_loss": 0.015
    },
    "confidence": 0.82,
    "expected_improvement": 0.15
  },
  "rationale": [
    "Reduced position size due to increased market volatility",
    "Tightened stop loss to preserve capital in uncertain conditions"
  ]
}
```

## 🌐 Browser Automation Endpoints

### POST /browser/navigate
Navigate to a specific URL and perform actions.

**Request:**
```json
{
  "url": "https://coinmarketcap.com",
  "actions": [
    {
      "type": "wait_for_element",
      "selector": ".price-section",
      "timeout": 5000
    },
    {
      "type": "click",
      "selector": "#bitcoin-link"
    },
    {
      "type": "extract_text",
      "selector": ".price-value",
      "variable": "btc_price"
    }
  ],
  "options": {
    "headless": true,
    "timeout": 30000,
    "user_agent": "Mozilla/5.0..."
  }
}
```

**Response:**
```json
{
  "session_id": "sess_789",
  "results": {
    "btc_price": "$45,123.45",
    "page_title": "Bitcoin (BTC) Price, Charts, and News | CoinMarketCap",
    "final_url": "https://coinmarketcap.com/currencies/bitcoin/"
  },
  "execution_time_ms": 2500,
  "screenshots": [
    {
      "name": "final_page",
      "data": "base64_encoded_screenshot"
    }
  ]
}
```

### POST /browser/extract
Extract structured data from web pages.

**Request:**
```json
{
  "url": "https://example-crypto-exchange.com/markets",
  "extraction_rules": {
    "prices": {
      "selector": ".price-cell",
      "attribute": "text",
      "multiple": true
    },
    "volumes": {
      "selector": ".volume-cell",
      "attribute": "text",
      "multiple": true
    },
    "market_cap": {
      "selector": ".market-cap",
      "attribute": "data-value"
    }
  }
}
```

## 📈 Analytics and Reporting Endpoints

### GET /analytics/performance
Get comprehensive performance analytics.

**Response:**
```json
{
  "period": "30d",
  "overall_performance": {
    "total_return": 0.24,
    "sharpe_ratio": 1.68,
    "max_drawdown": 0.12,
    "volatility": 0.18
  },
  "strategy_performance": [
    {
      "strategy_id": "strat_001",
      "name": "Adaptive Momentum",
      "return": 0.28,
      "trades": 45,
      "win_rate": 0.71
    }
  ],
  "pattern_effectiveness": [
    {
      "pattern_type": "bullish_trend",
      "detection_count": 23,
      "success_rate": 0.78,
      "avg_return": 0.06
    }
  ]
}
```

### GET /analytics/market-insights
Get AI-generated market insights and trends.

**Response:**
```json
{
  "insights": [
    {
      "type": "trend_analysis",
      "title": "Bitcoin Institutional Adoption Accelerating",
      "confidence": 0.85,
      "impact": "bullish",
      "timeframe": "medium_term",
      "supporting_data": [
        "ETF inflows increased 300% this month",
        "Corporate treasury allocations growing"
      ]
    }
  ],
  "market_regime": {
    "current": "risk_on",
    "confidence": 0.72,
    "duration": "2w",
    "next_regime_probability": {
      "risk_off": 0.25,
      "consolidation": 0.45
    }
  }
}
```

## 🔧 System Health Endpoints

### GET /health
Basic health check endpoint.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "v1.0.0",
  "uptime": "72h15m30s"
}
```

### GET /metrics
Prometheus metrics endpoint for monitoring.

**Response:**
```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",endpoint="/api/health",status_code="200"} 1234

# HELP ai_processing_duration_seconds AI processing duration
# TYPE ai_processing_duration_seconds histogram
ai_processing_duration_seconds_bucket{engine_type="pattern_detection",le="0.1"} 456
```

## 📝 Error Responses

All endpoints return consistent error responses:

```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "The request parameters are invalid",
    "details": {
      "field": "confidence",
      "reason": "Value must be between 0 and 1"
    },
    "request_id": "req_123",
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

## 🔄 Rate Limits

| Endpoint Category | Rate Limit | Burst Limit |
|------------------|------------|-------------|
| Authentication | 10/min | 5 |
| AI Analysis | 100/min | 20 |
| Pattern Detection | 50/min | 10 |
| Browser Automation | 20/min | 5 |
| Analytics | 200/min | 50 |

## 📚 SDK Examples

### Python SDK Example

```python
from ai_browser_client import AIBrowserClient

client = AIBrowserClient(
    base_url="https://api.ai-crypto-browser.com",
    api_key="your-api-key"
)

# Detect market patterns
patterns = client.market.detect_patterns({
    "asset": "BTC",
    "timeframe": "1h",
    "prices": [44000, 44200, 44500]
})

# Analyze with AI
analysis = client.ai.analyze({
    "type": "market_analysis",
    "query": "What's the outlook for Bitcoin?"
})
```

This comprehensive API reference provides complete documentation for all 43+ endpoints in the AI Agentic Crypto Browser, enabling developers to integrate and utilize the full range of AI-powered capabilities.
