# AI-Agentic Crypto Browser - API Documentation

## 🚀 Overview

The AI-Agentic Crypto Browser provides a comprehensive REST API for advanced AI-driven cryptocurrency analysis, trading, and browser automation. All endpoints require authentication unless otherwise specified.

## 🔐 Authentication

### Bearer Token Authentication
```http
Authorization: Bearer <your-jwt-token>
X-User-ID: <user-uuid>
```

### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "your-password"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "expires_at": "2024-01-01T12:00:00Z"
}
```

## 🧠 AI Agent Endpoints

### Enhanced AI Analysis
Comprehensive market analysis with multiple AI models.

```http
POST /ai/analyze
Content-Type: application/json
Authorization: Bearer <token>

{
  "symbols": ["BTC", "ETH", "ADA"],
  "timeframe": "1h",
  "indicators": ["RSI", "MACD", "Bollinger"],
  "analysis_types": ["technical", "sentiment", "risk", "prediction"]
}
```

**Response:**
```json
{
  "request_id": "req_123",
  "symbols": ["BTC", "ETH", "ADA"],
  "analysis": {
    "BTC": {
      "technical_analysis": {
        "trend": "bullish",
        "strength": 0.75,
        "indicators": {
          "RSI": 65.5,
          "MACD": 1.2,
          "Bollinger": "upper_band"
        }
      },
      "sentiment_analysis": {
        "score": 0.7,
        "label": "positive",
        "confidence": 0.85
      },
      "risk_assessment": {
        "overall_risk": 0.3,
        "risk_factors": [
          {
            "name": "volatility",
            "impact": "medium",
            "probability": "high",
            "score": 0.6
          }
        ]
      }
    }
  },
  "recommendation": "buy",
  "confidence": 0.82,
  "generated_at": "2024-01-01T12:00:00Z"
}
```

### Advanced Price Prediction
AI-powered price forecasting with multiple models.

```http
POST /ai/predict/price
Content-Type: application/json
Authorization: Bearer <token>

{
  "symbol": "BTC",
  "timeframes": ["1h", "4h", "1d", "1w"],
  "models": ["lstm", "transformer", "ensemble"],
  "features": ["price", "volume", "sentiment", "technical_indicators"]
}
```

**Response:**
```json
{
  "request_id": "pred_456",
  "symbol": "BTC",
  "predictions": {
    "1h": {
      "predicted_price": 52500.00,
      "confidence": 0.78,
      "price_range": {
        "min": 51800.00,
        "max": 53200.00
      },
      "probability_up": 0.65
    },
    "1d": {
      "predicted_price": 55000.00,
      "confidence": 0.72,
      "price_range": {
        "min": 52000.00,
        "max": 58000.00
      },
      "probability_up": 0.70
    }
  },
  "model_performance": {
    "lstm": {"accuracy": 0.75, "mse": 1250.5},
    "transformer": {"accuracy": 0.78, "mse": 1180.2},
    "ensemble": {"accuracy": 0.82, "mse": 1050.8}
  },
  "generated_at": "2024-01-01T12:00:00Z"
}
```

### Multi-Language Sentiment Analysis
Advanced sentiment analysis across multiple languages.

```http
POST /ai/analyze/sentiment
Content-Type: application/json
Authorization: Bearer <token>

{
  "texts": [
    "Bitcoin is showing incredible bullish momentum! 🚀",
    "El Bitcoin está mostrando un impulso alcista increíble"
  ],
  "languages": ["en", "es"],
  "options": {
    "detect_emotions": true,
    "analyze_intensity": true,
    "extract_aspects": true
  }
}
```

**Response:**
```json
{
  "request_id": "sent_789",
  "results": [
    {
      "text": "Bitcoin is showing incredible bullish momentum! 🚀",
      "language": "en",
      "sentiment": {
        "score": 0.85,
        "label": "positive",
        "confidence": 0.92,
        "intensity": 0.8,
        "emotions": {
          "joy": 0.7,
          "excitement": 0.8,
          "optimism": 0.9
        },
        "aspects": [
          {
            "aspect": "bitcoin",
            "sentiment": 0.85,
            "confidence": 0.9
          }
        ]
      }
    }
  ],
  "overall_sentiment": {
    "score": 0.82,
    "label": "positive",
    "confidence": 0.88
  },
  "generated_at": "2024-01-01T12:00:00Z"
}
```

### Comprehensive Predictive Analytics
Advanced predictive analytics with scenario modeling.

```http
POST /ai/analytics/predictive
Content-Type: application/json
Authorization: Bearer <token>

{
  "analysis_type": "comprehensive",
  "assets": ["BTC", "ETH"],
  "timeframe": "1d",
  "scenarios": ["bullish", "bearish", "sideways", "volatile"],
  "include_risk_analysis": true,
  "include_portfolio_impact": true,
  "include_correlations": true
}
```

**Response:**
```json
{
  "request_id": "analytics_101",
  "analysis_type": "comprehensive",
  "scenarios": [
    {
      "name": "bullish",
      "probability": 0.35,
      "expected_return": 0.15,
      "risk": 0.25,
      "duration": "7d",
      "triggers": ["institutional_adoption", "regulatory_clarity"]
    }
  ],
  "risk_metrics": {
    "overall_risk": 0.28,
    "var_95": 0.12,
    "max_drawdown": 0.18,
    "correlation_risk": 0.15
  },
  "portfolio_impact": {
    "expected_return": 0.12,
    "risk_adjusted_return": 0.08,
    "optimal_allocation": {
      "BTC": 0.6,
      "ETH": 0.4
    }
  },
  "generated_at": "2024-01-01T12:00:00Z"
}
```

### Advanced NLP Analysis
Comprehensive natural language processing.

```http
POST /ai/nlp/analyze
Content-Type: application/json
Authorization: Bearer <token>

{
  "texts": [
    "Bitcoin's institutional adoption is accelerating with major corporations adding BTC to their treasury reserves."
  ],
  "sources": ["news"],
  "options": {
    "detect_language": true,
    "extract_entities": true,
    "perform_topic_modeling": true,
    "classify_text": true,
    "analyze_sentiment": true,
    "extract_keywords": true,
    "detect_emotions": true,
    "analyze_readability": true
  }
}
```

**Response:**
```json
{
  "request_id": "nlp_202",
  "results": [
    {
      "text": "Bitcoin's institutional adoption...",
      "original_language": "en",
      "sentiment": {
        "score": 0.75,
        "label": "positive",
        "confidence": 0.88
      },
      "entities": [
        {
          "text": "Bitcoin",
          "type": "CRYPTO",
          "confidence": 0.95,
          "start_pos": 0,
          "end_pos": 7
        }
      ],
      "topics": [
        {
          "id": "institutional_adoption",
          "name": "Institutional Adoption",
          "probability": 0.85,
          "keywords": ["institutional", "adoption", "corporations"]
        }
      ],
      "classification": {
        "category": "financial",
        "confidence": 0.92
      },
      "keywords": [
        {
          "keyword": "institutional",
          "score": 0.8,
          "frequency": 2
        }
      ],
      "readability_score": 75.5
    }
  ],
  "aggregated_results": {
    "overall_sentiment": {
      "score": 0.75,
      "label": "positive"
    },
    "top_entities": [...],
    "top_topics": [...],
    "language_distribution": {"en": 1.0}
  },
  "generated_at": "2024-01-01T12:00:00Z"
}
```

### Intelligent Decision Making
AI-driven trading decisions with risk management.

```http
POST /ai/decisions/request
Content-Type: application/json
Authorization: Bearer <token>

{
  "decision_type": "trade",
  "context": {
    "market_conditions": "bullish",
    "time_horizon": "short",
    "urgency": "medium",
    "trigger_event": "price_breakout",
    "technical_indicators": {
      "rsi": 30.0,
      "macd_signal": 1.0,
      "volume_ratio": 1.5
    }
  },
  "constraints": {
    "max_position_size": 1000.0,
    "max_risk_exposure": 0.05,
    "allowed_assets": ["BTC", "ETH"]
  },
  "preferences": {
    "risk_tolerance": 0.6,
    "auto_execution_level": "none",
    "decision_speed": "normal"
  },
  "options": {
    "require_confirmation": true,
    "explain_reasoning": true,
    "simulate_execution": true
  }
}
```

**Response:**
```json
{
  "request_id": "req_303",
  "decision_id": "dec_404",
  "decision_type": "trade",
  "recommendation": {
    "action": "buy",
    "asset": "BTC",
    "quantity": 0.02,
    "order_type": "market",
    "confidence": 0.82,
    "expected_return": 0.08,
    "risk_score": 0.25,
    "reasoning": "Strong bullish momentum with favorable risk-reward ratio"
  },
  "alternatives": [
    {
      "action": "buy",
      "asset": "ETH",
      "quantity": 0.3,
      "confidence": 0.75
    }
  ],
  "risk_assessment": {
    "overall_risk": 0.28,
    "risk_factors": [
      {
        "risk": "market_volatility",
        "probability": 0.4,
        "impact": 0.2,
        "severity": "medium"
      }
    ],
    "max_loss": 50.0,
    "confidence": 0.85
  },
  "reasoning": {
    "primary_factors": [
      {
        "factor": "technical_breakout",
        "weight": 0.4,
        "impact": "positive",
        "confidence": 0.8
      }
    ],
    "confidence": 0.82
  },
  "expected_outcome": {
    "expected_return": 0.08,
    "expected_risk": 0.25,
    "probability_of_profit": 0.75,
    "time_horizon": "24h"
  },
  "execution_plan": {
    "steps": [
      {
        "step_id": "step_1",
        "action": "place_order",
        "estimated_time": "30s",
        "parameters": {
          "asset": "BTC",
          "quantity": 0.02,
          "type": "market"
        }
      }
    ],
    "total_estimated_time": "30s",
    "estimated_cost": 10.0
  },
  "requires_approval": true,
  "auto_executable": false,
  "generated_at": "2024-01-01T12:00:00Z",
  "expires_at": "2024-01-02T12:00:00Z"
}
```

## 📚 Learning and Adaptation Endpoints

### User Behavior Learning
Learn from user trading behavior and preferences.

```http
POST /ai/learning/behavior
Content-Type: application/json
Authorization: Bearer <token>

{
  "type": "trade",
  "data": {
    "symbol": "BTC",
    "side": "buy",
    "amount": 1000.0,
    "risk_level": 0.7,
    "decision_factors": {
      "technical_analysis": 0.6,
      "sentiment_analysis": 0.4
    }
  },
  "context": {
    "market_condition": "bullish",
    "volatility": 0.3
  },
  "outcome": "success",
  "performance": 0.05
}
```

### Get User Profile
Retrieve learned user profile and preferences.

```http
GET /ai/learning/profile
Authorization: Bearer <token>
```

**Response:**
```json
{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "risk_tolerance": 0.65,
  "trading_style": "moderate",
  "behavior_score": 0.75,
  "trading_patterns": {
    "avg_holding_period": "2d",
    "trading_frequency": 3.5,
    "preferred_assets": ["BTC", "ETH"],
    "risk_management": {
      "stop_loss_usage": 0.8,
      "position_sizing": "kelly_criterion"
    }
  },
  "decision_factors": {
    "technical_analysis": 0.6,
    "sentiment_analysis": 0.3,
    "fundamental_analysis": 0.1
  },
  "performance_metrics": {
    "total_return": 0.15,
    "win_rate": 0.68,
    "sharpe_ratio": 1.2,
    "max_drawdown": 0.08
  },
  "last_updated": "2024-01-01T12:00:00Z"
}
```

### Get Market Patterns
Retrieve discovered market patterns and insights.

```http
GET /ai/learning/patterns
Authorization: Bearer <token>
```

### Request Model Adaptation
Request adaptation for AI models based on performance.

```http
POST /ai/adaptation/request
Content-Type: application/json
Authorization: Bearer <token>

{
  "model_id": "price_prediction",
  "type": "performance",
  "trigger": "accuracy_decline",
  "priority": 2,
  "data": {
    "performance_drop": 0.1,
    "sample_size": 1000
  }
}
```

## 📊 Performance and Monitoring Endpoints

### Get Decision History
Retrieve historical decisions and their outcomes.

```http
GET /ai/decisions/history?limit=50
Authorization: Bearer <token>
```

### Get Performance Metrics
Retrieve AI system performance metrics.

```http
GET /ai/decisions/performance
Authorization: Bearer <token>
```

### Get Active Decisions
Retrieve currently active/processing decisions.

```http
GET /ai/decisions/active
Authorization: Bearer <token>
```

## 🎨 Multi-Modal AI Endpoints

### Comprehensive Multi-Modal Analysis
Process multiple types of content (images, documents, audio) in a single request.

```http
POST /ai/multimodal/analyze
Content-Type: application/json
Authorization: Bearer <token>

{
  "type": "mixed",
  "content": [
    {
      "id": "content-1",
      "type": "image",
      "data": "base64-encoded-image-data",
      "mime_type": "image/png",
      "filename": "chart.png",
      "size": 12345
    },
    {
      "id": "content-2",
      "type": "document",
      "data": "base64-encoded-document-data",
      "mime_type": "application/pdf",
      "filename": "report.pdf",
      "size": 54321
    }
  ],
  "options": {
    "analyze_images": true,
    "extract_text": true,
    "analyze_charts": true,
    "analyze_sentiment": true,
    "generate_summary": true
  }
}
```

**Response:**
```json
{
  "request_id": "req_123",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "type": "mixed",
  "results": [
    {
      "content_id": "content-1",
      "type": "image",
      "image_analysis": {
        "objects": [
          {
            "label": "chart",
            "confidence": 0.95,
            "bounding_box": {
              "x": 0.1,
              "y": 0.1,
              "width": 0.8,
              "height": 0.8
            }
          }
        ],
        "trading_signals": [
          {
            "type": "support",
            "signal": "buy",
            "confidence": 0.8,
            "price": 50000.0,
            "description": "Strong support level identified"
          }
        ]
      },
      "chart_analysis": {
        "chart_type": "candlestick",
        "asset": "BTC",
        "technical_signals": [
          {
            "indicator": "RSI",
            "signal": "buy",
            "value": 30.0,
            "confidence": 0.85
          }
        ],
        "recommendation": {
          "action": "buy",
          "confidence": 0.8,
          "target": 55000.0,
          "stop_loss": 48000.0
        }
      },
      "confidence": 0.85
    }
  ],
  "aggregated_data": {
    "key_insights": [
      "Strong bullish momentum detected",
      "Multiple buy signals confirmed"
    ],
    "trading_signals": [...],
    "summary": "Analysis shows strong bullish indicators across multiple data sources",
    "processing_stats": {
      "total_items": 2,
      "successful_items": 2,
      "processing_time": "1.5s"
    }
  },
  "processing_time": "1.5s",
  "generated_at": "2024-01-01T12:00:00Z"
}
```

### Image Analysis
Upload and analyze images for trading charts, patterns, and signals.

```http
POST /ai/multimodal/image
Content-Type: multipart/form-data
Authorization: Bearer <token>

Form Data:
- image: [image file]
- extract_text: true
- analyze_charts: true
- detect_objects: true
- generate_summary: true
```

### Document Analysis
Upload and analyze documents for financial insights and trading information.

```http
POST /ai/multimodal/document
Content-Type: multipart/form-data
Authorization: Bearer <token>

Form Data:
- document: [document file]
- analyze_sentiment: true
- extract_entities: true
- generate_summary: true
```

### Audio Analysis
Process audio files for voice commands and trading instructions.

```http
POST /ai/multimodal/audio
Content-Type: multipart/form-data
Authorization: Bearer <token>

Form Data:
- audio: [audio file]
- analyze_sentiment: true
- extract_entities: true
```

### Chart Analysis
Specialized analysis for trading charts and financial visualizations.

```http
POST /ai/multimodal/chart
Content-Type: multipart/form-data
Authorization: Bearer <token>

Form Data:
- chart: [chart image file]
```

### Get Supported Formats
Retrieve supported file formats for multi-modal analysis.

```http
GET /ai/multimodal/formats
Authorization: Bearer <token>
```

**Response:**
```json
{
  "supported_formats": {
    "images": ["jpg", "jpeg", "png", "gif", "webp"],
    "documents": ["pdf", "docx", "txt", "csv", "xlsx"],
    "audio": ["mp3", "wav", "m4a", "ogg", "flac"]
  },
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## 🧠 User Behavior Learning Endpoints

### Learn from User Behavior
Record and learn from user behavior events to build personalized profiles.

```http
POST /ai/behavior/learn
Content-Type: application/json
Authorization: Bearer <token>

{
  "type": "trade",
  "action": "buy_btc",
  "context": {
    "market_conditions": "bullish",
    "portfolio_state": {
      "btc_balance": 0.5
    },
    "time_of_day": "morning",
    "day_of_week": "monday",
    "session_duration": "2h",
    "previous_actions": ["analyze_chart", "check_news"],
    "emotional_state": "confident",
    "information_sources": ["technical_analysis", "news"],
    "external_factors": {
      "market_sentiment": "positive"
    }
  },
  "outcome": {
    "success": true,
    "performance": 0.05,
    "satisfaction": 0.8,
    "time_to_decision": "15m",
    "confidence_level": 0.7,
    "regret": 0.1,
    "learning_value": 0.8
  },
  "duration": "30m"
}
```

**Response:**
```json
{
  "success": true,
  "event_id": "evt_123",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Get User Behavior Profile
Retrieve comprehensive user behavior profile with trading style, risk profile, personality, and preferences.

```http
GET /ai/behavior/profile
Authorization: Bearer <token>
```

**Response:**
```json
{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2024-01-01T10:00:00Z",
  "last_updated": "2024-01-01T12:00:00Z",
  "observation_count": 25,
  "confidence": 0.78,
  "trading_style": {
    "primary_style": "day_trader",
    "secondary_styles": ["swing_trader"],
    "trading_frequency": 3.5,
    "average_hold_time": "6h",
    "preferred_timeframes": ["1h", "4h"],
    "preferred_assets": ["BTC", "ETH"],
    "decision_speed": "medium",
    "analysis_depth": "moderate",
    "confidence": 0.85
  },
  "risk_profile": {
    "risk_tolerance": 0.7,
    "risk_capacity": 0.8,
    "max_drawdown_tolerance": 0.15,
    "position_sizing_style": "percentage",
    "stop_loss_usage": 0.8,
    "take_profit_usage": 0.6,
    "leverage_comfort": 0.3,
    "volatility_tolerance": 0.5,
    "emotional_stability": 0.7,
    "loss_aversion": 0.6,
    "confidence": 0.75
  },
  "personality_profile": {
    "trader_type": "analytical",
    "decision_making": "rational",
    "information_processing": "sequential",
    "planning_orientation": "structured",
    "stress_response": "calm",
    "confidence_level": "moderate",
    "learning_style": "visual",
    "social_influence": 0.3,
    "patience": 0.8,
    "discipline": 0.9,
    "adaptability": 0.7,
    "traits": {
      "openness": 0.7,
      "conscientiousness": 0.8,
      "extraversion": 0.5,
      "agreeableness": 0.6,
      "neuroticism": 0.3
    },
    "biases": [
      {
        "type": "confirmation",
        "strength": 0.4,
        "frequency": 0.3,
        "impact": 0.2
      }
    ],
    "confidence": 0.72
  },
  "performance_metrics": {
    "overall_return": 0.15,
    "annualized_return": 0.18,
    "volatility": 0.25,
    "sharpe_ratio": 0.72,
    "max_drawdown": 0.08,
    "win_rate": 0.65,
    "profit_factor": 1.8,
    "total_trades": 45,
    "successful_trades": 29,
    "confidence": 0.80
  },
  "learning_progress": {
    "total_observations": 25,
    "learning_rate": 0.1,
    "model_accuracy": 0.78,
    "prediction_accuracy": 0.72,
    "pattern_recognition": 0.68,
    "preference_stability": 0.85,
    "profile_completeness": 0.75,
    "learning_velocity": 2.5,
    "adaptation_rate": 0.15,
    "confidence_growth": 0.12,
    "milestones": [
      {
        "id": "first_10_observations",
        "name": "First 10 Observations",
        "description": "Completed first 10 behavior observations",
        "achieved_at": "2024-01-01T11:00:00Z",
        "value": 10,
        "threshold": 10
      }
    ]
  }
}
```

### Get Personalized Recommendations
Retrieve AI-generated personalized recommendations based on user behavior profile.

```http
GET /ai/behavior/recommendations?limit=5
Authorization: Bearer <token>
```

**Response:**
```json
{
  "recommendations": [
    {
      "id": "rec_123",
      "type": "strategy",
      "title": "Optimize Your Position Sizing",
      "description": "Based on your trading patterns, consider adjusting your position sizing strategy to improve risk-adjusted returns",
      "reasoning": [
        "Your win rate is above average at 65%",
        "Risk-adjusted returns could be improved with better position sizing",
        "Your emotional stability suggests you can handle larger positions"
      ],
      "confidence": 0.82,
      "priority": "high",
      "category": "optimization",
      "action_required": false,
      "expected_outcome": {
        "probability_of_success": 0.75,
        "expected_return": 0.08,
        "expected_risk": 0.03,
        "time_horizon": "30d"
      },
      "risk_assessment": {
        "overall_risk": 0.25,
        "risk_factors": [
          {
            "type": "market_risk",
            "description": "Market volatility could impact results",
            "impact": 0.3,
            "probability": 0.4,
            "severity": "medium"
          }
        ],
        "mitigation": [
          "Start with small position size increases",
          "Monitor performance closely",
          "Set strict stop losses"
        ],
        "max_loss": 0.02,
        "risk_reward": 2.7
      },
      "personalization": {
        "personalization_score": 0.85,
        "user_factors": ["trading_style", "risk_tolerance", "performance_history"],
        "behavior_factors": ["win_rate", "position_sizing", "emotional_stability"],
        "context_factors": ["market_conditions", "portfolio_state"],
        "adaptations": ["personalized_thresholds", "risk_adjusted_targets"]
      },
      "created_at": "2024-01-01T12:00:00Z",
      "status": "pending"
    }
  ],
  "count": 1,
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Get Behavior History
Retrieve user's behavior event history for analysis and review.

```http
GET /ai/behavior/history?limit=10
Authorization: Bearer <token>
```

**Response:**
```json
{
  "history": [
    {
      "id": "evt_123",
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "type": "trade",
      "action": "buy_btc",
      "context": {
        "market_conditions": "bullish",
        "time_of_day": "morning",
        "session_duration": "2h"
      },
      "outcome": {
        "success": true,
        "performance": 0.05,
        "satisfaction": 0.8
      },
      "timestamp": "2024-01-01T12:00:00Z",
      "duration": "30m"
    }
  ],
  "count": 1,
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Update Recommendation Status
Update the status of a personalized recommendation (accept, reject, etc.).

```http
PUT /ai/behavior/recommendation/{recommendation_id}/status
Content-Type: application/json
Authorization: Bearer <token>

{
  "status": "accepted"
}
```

**Valid statuses:** `pending`, `accepted`, `rejected`, `expired`

**Response:**
```json
{
  "success": true,
  "recommendation_id": "rec_123",
  "status": "accepted",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Get Learning Models
Retrieve information about the machine learning models used for behavior analysis.

```http
GET /ai/behavior/models
Authorization: Bearer <token>
```

**Response:**
```json
{
  "models": {
    "preference_model": {
      "id": "pref_model_v1",
      "type": "neural_network",
      "purpose": "preference",
      "version": "1.0.0",
      "accuracy": 0.85,
      "training_data": 1000,
      "last_trained": "2024-01-01T10:00:00Z",
      "next_training": "2024-01-02T10:00:00Z"
    },
    "pattern_model": {
      "id": "pattern_model_v1",
      "type": "ensemble",
      "purpose": "pattern",
      "version": "1.0.0",
      "accuracy": 0.78,
      "training_data": 2500,
      "last_trained": "2024-01-01T10:00:00Z",
      "next_training": "2024-01-02T10:00:00Z"
    }
  },
  "count": 2,
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## 🌐 Browser Automation Endpoints

### Create Browser Session
```http
POST /browser/sessions
Content-Type: application/json
Authorization: Bearer <token>

{
  "headless": true,
  "viewport": {"width": 1920, "height": 1080},
  "user_agent": "custom-agent"
}
```

### Navigate to URL
```http
POST /browser/navigate
Content-Type: application/json
Authorization: Bearer <token>
X-Session-ID: <session-id>

{
  "url": "https://example.com",
  "wait_for": "networkidle"
}
```

## 🔗 Web3 Integration Endpoints

### Connect Wallet
```http
POST /web3/connect-wallet
Content-Type: application/json
Authorization: Bearer <token>

{
  "wallet_type": "metamask",
  "address": "0x742d35Cc6634C0532925a3b8D4C9db96C4b4Db45"
}
```

### Get Balance
```http
GET /web3/balance?address=0x742d35Cc6634C0532925a3b8D4C9db96C4b4Db45
Authorization: Bearer <token>
```

## 📋 Error Handling

All endpoints return consistent error responses:

```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "The request parameters are invalid",
    "details": {
      "field": "symbol",
      "reason": "Symbol 'INVALID' is not supported"
    }
  },
  "request_id": "req_123",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Common Error Codes
- `UNAUTHORIZED` (401) - Invalid or missing authentication
- `FORBIDDEN` (403) - Insufficient permissions
- `NOT_FOUND` (404) - Resource not found
- `INVALID_REQUEST` (400) - Invalid request parameters
- `RATE_LIMITED` (429) - Too many requests
- `INTERNAL_ERROR` (500) - Server error

## 🚀 Rate Limiting

API endpoints are rate limited to ensure fair usage:

- **Standard endpoints:** 100 requests/minute
- **AI analysis endpoints:** 20 requests/minute
- **Decision making endpoints:** 10 requests/minute

Rate limit headers are included in responses:
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

This comprehensive API enables sophisticated AI-driven cryptocurrency analysis, trading, and automation capabilities through a clean, RESTful interface.
