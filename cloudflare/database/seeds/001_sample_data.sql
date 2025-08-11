-- Seed data for testing AI Agentic Crypto Browser

-- Sample users
INSERT INTO users (id, email, password_hash, name, role, created_at) VALUES
('user-1', 'demo@example.com', 'hashed_password_demo', 'Demo User', 'user', '2024-01-01T00:00:00Z'),
('user-2', 'admin@example.com', 'hashed_password_admin', 'Admin User', 'admin', '2024-01-01T00:00:00Z'),
('user-3', 'trader@example.com', 'hashed_password_trader', 'Pro Trader', 'user', '2024-01-01T00:00:00Z');

-- Sample market data
INSERT INTO market_data (symbol, price, volume_24h, change_24h, market_cap, last_updated) VALUES
('BTCUSDT', 45500.00, 28500000000, 2.5, 890000000000, '2024-01-01T12:00:00Z'),
('ETHUSDT', 2800.00, 15200000000, -1.2, 336000000000, '2024-01-01T12:00:00Z'),
('ADAUSDT', 0.45, 850000000, 5.8, 15800000000, '2024-01-01T12:00:00Z'),
('SOLUSDT', 98.50, 2100000000, 3.2, 43500000000, '2024-01-01T12:00:00Z'),
('DOTUSDT', 7.25, 420000000, -0.8, 9200000000, '2024-01-01T12:00:00Z');

-- Sample conversations
INSERT INTO conversations (id, user_id, title, created_at) VALUES
('conv-1', 'user-1', 'Bitcoin Analysis Discussion', '2024-01-01T10:00:00Z'),
('conv-2', 'user-1', 'Portfolio Optimization Help', '2024-01-01T11:00:00Z'),
('conv-3', 'user-3', 'Trading Strategy Review', '2024-01-01T09:00:00Z');

-- Sample messages
INSERT INTO messages (id, conversation_id, role, content, created_at, model_used) VALUES
('msg-1', 'conv-1', 'user', 'What do you think about Bitcoin''s current price action?', '2024-01-01T10:00:00Z', NULL),
('msg-2', 'conv-1', 'assistant', 'Based on current technical indicators, Bitcoin is showing bullish momentum with strong support at $44,000. The RSI is at 65.5, indicating healthy upward movement without being overbought.', '2024-01-01T10:00:30Z', 'gpt-3.5-turbo'),
('msg-3', 'conv-2', 'user', 'How should I diversify my crypto portfolio?', '2024-01-01T11:00:00Z', NULL),
('msg-4', 'conv-2', 'assistant', 'For optimal diversification, consider allocating 40% to Bitcoin, 25% to Ethereum, 15% to large-cap altcoins, 15% to DeFi tokens, and 5% to emerging projects. This balances stability with growth potential.', '2024-01-01T11:00:30Z', 'gpt-3.5-turbo');

-- Sample portfolio holdings
INSERT INTO portfolio_holdings (id, user_id, symbol, quantity, average_cost, current_price, exchange) VALUES
('holding-1', 'user-1', 'BTC', 0.25, 42000.00, 45500.00, 'binance'),
('holding-2', 'user-1', 'ETH', 3.2, 2500.00, 2800.00, 'binance'),
('holding-3', 'user-1', 'ADA', 1000.0, 0.40, 0.45, 'binance'),
('holding-4', 'user-3', 'BTC', 1.5, 40000.00, 45500.00, 'coinbase'),
('holding-5', 'user-3', 'SOL', 50.0, 85.00, 98.50, 'coinbase');

-- Sample trading orders
INSERT INTO trading_orders (id, user_id, symbol, side, type, quantity, price, status, created_at, exchange) VALUES
('order-1', 'user-1', 'BTCUSDT', 'buy', 'limit', 0.1, 44000.00, 'pending', '2024-01-01T12:30:00Z', 'binance'),
('order-2', 'user-3', 'ETHUSDT', 'sell', 'limit', 1.0, 2850.00, 'pending', '2024-01-01T12:45:00Z', 'coinbase'),
('order-3', 'user-1', 'ADAUSDT', 'buy', 'market', 500.0, NULL, 'filled', '2024-01-01T11:15:00Z', 'binance');

-- Sample AI analysis
INSERT INTO ai_analysis (id, user_id, symbol, analysis_type, timeframe, result, confidence, created_at) VALUES
('analysis-1', 'user-1', 'BTCUSDT', 'technical', '1h', '{"trend": "bullish", "rsi": 65.5, "macd": "bullish_crossover", "recommendation": "HOLD"}', 0.85, '2024-01-01T12:00:00Z'),
('analysis-2', 'user-3', 'ETHUSDT', 'sentiment', '4h', '{"sentiment": "positive", "social_score": 7.2, "news_sentiment": "bullish", "recommendation": "BUY"}', 0.78, '2024-01-01T11:30:00Z');

-- Sample price predictions
INSERT INTO price_predictions (id, symbol, predicted_price, current_price, timeframe, confidence, model_used, created_at, target_time) VALUES
('pred-1', 'BTCUSDT', 47200.00, 45500.00, '24h', 0.72, 'ensemble_v1', '2024-01-01T12:00:00Z', '2024-01-02T12:00:00Z'),
('pred-2', 'ETHUSDT', 2950.00, 2800.00, '24h', 0.68, 'ensemble_v1', '2024-01-01T12:00:00Z', '2024-01-02T12:00:00Z');

-- Sample risk assessments
INSERT INTO risk_assessments (id, user_id, risk_score, risk_level, factors, recommendations, portfolio_value, created_at) VALUES
('risk-1', 'user-1', 6.5, 'moderate', '["Market volatility", "Portfolio concentration"]', '["Diversify holdings", "Set stop losses"]', 15000.00, '2024-01-01T12:00:00Z'),
('risk-2', 'user-3', 8.2, 'high', '["High leverage", "Concentrated positions"]', '["Reduce leverage", "Increase diversification", "Implement risk management"]', 75000.00, '2024-01-01T12:00:00Z');

-- Sample user preferences
INSERT INTO user_preferences (user_id, theme, language, timezone, notifications, trading_preferences, created_at) VALUES
('user-1', 'dark', 'en', 'UTC', '{"email": true, "push": false, "sms": false}', '{"auto_trading": false, "risk_level": "moderate"}', '2024-01-01T00:00:00Z'),
('user-3', 'light', 'en', 'EST', '{"email": true, "push": true, "sms": true}', '{"auto_trading": true, "risk_level": "high"}', '2024-01-01T00:00:00Z');
