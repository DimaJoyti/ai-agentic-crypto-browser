import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const responseTime = new Trend('response_time');
const tradingRequests = new Counter('trading_requests');
const authRequests = new Counter('auth_requests');

// Test configuration
export const options = {
  stages: [
    // Ramp-up
    { duration: '2m', target: 10 },   // Ramp up to 10 users
    { duration: '5m', target: 50 },   // Ramp up to 50 users
    { duration: '10m', target: 100 }, // Ramp up to 100 users
    
    // Steady state
    { duration: '15m', target: 100 }, // Stay at 100 users
    
    // Peak load
    { duration: '5m', target: 200 },  // Spike to 200 users
    { duration: '10m', target: 200 }, // Stay at 200 users
    
    // Ramp-down
    { duration: '5m', target: 50 },   // Ramp down to 50 users
    { duration: '2m', target: 0 },    // Ramp down to 0 users
  ],
  thresholds: {
    // Performance thresholds
    http_req_duration: ['p(95)<500'], // 95% of requests must complete below 500ms
    http_req_failed: ['rate<0.1'],    // Error rate must be below 10%
    
    // Custom thresholds
    errors: ['rate<0.05'],            // Error rate must be below 5%
    response_time: ['p(99)<1000'],    // 99% of requests must complete below 1s
  },
};

// Test data
const BASE_URL = __ENV.STAGING_URL || 'http://localhost:8080';
const FRONTEND_URL = __ENV.FRONTEND_URL || 'http://localhost:3000';

const testUsers = [
  { username: 'trader1@example.com', password: 'password123' },
  { username: 'trader2@example.com', password: 'password123' },
  { username: 'trader3@example.com', password: 'password123' },
];

const tradingPairs = ['BTC/USD', 'ETH/USD', 'BNB/USD', 'ADA/USD', 'DOT/USD'];
const orderTypes = ['market', 'limit', 'stop', 'stop-limit'];

// Utility functions
function getRandomUser() {
  return testUsers[Math.floor(Math.random() * testUsers.length)];
}

function getRandomTradingPair() {
  return tradingPairs[Math.floor(Math.random() * tradingPairs.length)];
}

function getRandomOrderType() {
  return orderTypes[Math.floor(Math.random() * orderTypes.length)];
}

function generateRandomPrice(base = 50000) {
  return base + (Math.random() - 0.5) * base * 0.1;
}

function generateRandomQuantity() {
  return Math.random() * 10 + 0.1;
}

// Authentication
function authenticate() {
  const user = getRandomUser();
  
  const loginPayload = {
    email: user.username,
    password: user.password,
  };

  const loginResponse = http.post(`${BASE_URL}/api/auth/login`, JSON.stringify(loginPayload), {
    headers: { 'Content-Type': 'application/json' },
  });

  authRequests.add(1);

  const loginSuccess = check(loginResponse, {
    'login status is 200': (r) => r.status === 200,
    'login response has token': (r) => r.json('token') !== undefined,
  });

  if (!loginSuccess) {
    errorRate.add(1);
    return null;
  }

  return loginResponse.json('token');
}

// API Tests
export default function () {
  // Authenticate user
  const token = authenticate();
  if (!token) {
    return;
  }

  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
  };

  // Test scenarios with different weights
  const scenario = Math.random();

  if (scenario < 0.3) {
    // 30% - Market data requests
    testMarketData(headers);
  } else if (scenario < 0.5) {
    // 20% - Portfolio operations
    testPortfolio(headers);
  } else if (scenario < 0.7) {
    // 20% - Trading operations
    testTrading(headers);
  } else if (scenario < 0.85) {
    // 15% - User dashboard
    testDashboard(headers);
  } else {
    // 15% - Frontend pages
    testFrontend();
  }

  sleep(Math.random() * 3 + 1); // Random sleep between 1-4 seconds
}

function testMarketData(headers) {
  const pair = getRandomTradingPair();
  
  // Get market data
  const marketDataResponse = http.get(`${BASE_URL}/api/market/data/${pair}`, { headers });
  
  const success = check(marketDataResponse, {
    'market data status is 200': (r) => r.status === 200,
    'market data has price': (r) => r.json('price') !== undefined,
    'market data response time < 200ms': (r) => r.timings.duration < 200,
  });

  responseTime.add(marketDataResponse.timings.duration);
  if (!success) errorRate.add(1);

  // Get order book
  const orderBookResponse = http.get(`${BASE_URL}/api/market/orderbook/${pair}`, { headers });
  
  check(orderBookResponse, {
    'order book status is 200': (r) => r.status === 200,
    'order book has bids': (r) => r.json('bids') !== undefined,
    'order book has asks': (r) => r.json('asks') !== undefined,
  });

  responseTime.add(orderBookResponse.timings.duration);
}

function testPortfolio(headers) {
  // Get portfolio
  const portfolioResponse = http.get(`${BASE_URL}/api/portfolio`, { headers });
  
  const success = check(portfolioResponse, {
    'portfolio status is 200': (r) => r.status === 200,
    'portfolio has balance': (r) => r.json('totalBalance') !== undefined,
    'portfolio response time < 300ms': (r) => r.timings.duration < 300,
  });

  responseTime.add(portfolioResponse.timings.duration);
  if (!success) errorRate.add(1);

  // Get portfolio history
  const historyResponse = http.get(`${BASE_URL}/api/portfolio/history?period=7d`, { headers });
  
  check(historyResponse, {
    'portfolio history status is 200': (r) => r.status === 200,
    'portfolio history has data': (r) => r.json('data') !== undefined,
  });

  responseTime.add(historyResponse.timings.duration);
}

function testTrading(headers) {
  const pair = getRandomTradingPair();
  const orderType = getRandomOrderType();
  
  // Create order
  const orderPayload = {
    symbol: pair,
    side: Math.random() > 0.5 ? 'buy' : 'sell',
    type: orderType,
    quantity: generateRandomQuantity(),
    price: orderType !== 'market' ? generateRandomPrice() : undefined,
    timeInForce: 'GTC',
  };

  const createOrderResponse = http.post(
    `${BASE_URL}/api/trading/orders`,
    JSON.stringify(orderPayload),
    { headers }
  );

  tradingRequests.add(1);

  const orderSuccess = check(createOrderResponse, {
    'create order status is 200 or 201': (r) => [200, 201].includes(r.status),
    'create order has orderId': (r) => r.json('orderId') !== undefined,
    'create order response time < 500ms': (r) => r.timings.duration < 500,
  });

  responseTime.add(createOrderResponse.timings.duration);
  if (!orderSuccess) errorRate.add(1);

  // Get orders
  const ordersResponse = http.get(`${BASE_URL}/api/trading/orders`, { headers });
  
  check(ordersResponse, {
    'get orders status is 200': (r) => r.status === 200,
    'get orders has data': (r) => Array.isArray(r.json()),
  });

  responseTime.add(ordersResponse.timings.duration);

  // Get trading history
  const historyResponse = http.get(`${BASE_URL}/api/trading/history?limit=50`, { headers });
  
  check(historyResponse, {
    'trading history status is 200': (r) => r.status === 200,
    'trading history response time < 400ms': (r) => r.timings.duration < 400,
  });

  responseTime.add(historyResponse.timings.duration);
}

function testDashboard(headers) {
  // Get dashboard data
  const dashboardResponse = http.get(`${BASE_URL}/api/dashboard`, { headers });
  
  const success = check(dashboardResponse, {
    'dashboard status is 200': (r) => r.status === 200,
    'dashboard has summary': (r) => r.json('summary') !== undefined,
    'dashboard response time < 600ms': (r) => r.timings.duration < 600,
  });

  responseTime.add(dashboardResponse.timings.duration);
  if (!success) errorRate.add(1);

  // Get notifications
  const notificationsResponse = http.get(`${BASE_URL}/api/notifications`, { headers });
  
  check(notificationsResponse, {
    'notifications status is 200': (r) => r.status === 200,
    'notifications is array': (r) => Array.isArray(r.json()),
  });

  responseTime.add(notificationsResponse.timings.duration);

  // Get user preferences
  const preferencesResponse = http.get(`${BASE_URL}/api/user/preferences`, { headers });
  
  check(preferencesResponse, {
    'preferences status is 200': (r) => r.status === 200,
    'preferences has theme': (r) => r.json('theme') !== undefined,
  });

  responseTime.add(preferencesResponse.timings.duration);
}

function testFrontend() {
  // Test frontend pages
  const pages = ['/', '/trading', '/portfolio', '/analytics', '/ux-demo'];
  const page = pages[Math.floor(Math.random() * pages.length)];
  
  const pageResponse = http.get(`${FRONTEND_URL}${page}`);
  
  const success = check(pageResponse, {
    'frontend page status is 200': (r) => r.status === 200,
    'frontend page has content': (r) => r.body.length > 1000,
    'frontend page response time < 1000ms': (r) => r.timings.duration < 1000,
  });

  responseTime.add(pageResponse.timings.duration);
  if (!success) errorRate.add(1);

  // Test static assets
  const assetResponse = http.get(`${FRONTEND_URL}/_next/static/css/app.css`);
  
  check(assetResponse, {
    'static asset loads': (r) => r.status === 200,
    'static asset response time < 200ms': (r) => r.timings.duration < 200,
  });

  responseTime.add(assetResponse.timings.duration);
}

// Teardown function
export function teardown(data) {
  console.log('Load test completed');
  console.log(`Total trading requests: ${tradingRequests.count}`);
  console.log(`Total auth requests: ${authRequests.count}`);
  console.log(`Error rate: ${(errorRate.rate * 100).toFixed(2)}%`);
  console.log(`Average response time: ${responseTime.avg.toFixed(2)}ms`);
  console.log(`95th percentile response time: ${responseTime.p(95).toFixed(2)}ms`);
}
