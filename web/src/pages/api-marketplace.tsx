import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { 
  Code, 
  Key, 
  DollarSign, 
  TrendingUp, 
  Shield, 
  Zap,
  Brain,
  BarChart3,
  Copy,
  Eye,
  EyeOff
} from 'lucide-react';

interface APIEndpoint {
  id: string;
  endpoint: string;
  price_per_call: string;
  category: string;
  description: string;
  rate_limit: number;
  requires_plan: boolean;
}

interface APIKey {
  key_id: string;
  name: string;
  description: string;
  scopes: string[];
  created_at: string;
  last_used: string;
  active: boolean;
}

const APIMarketplace = () => {
  const [endpoints, setEndpoints] = useState<Record<string, APIEndpoint>>({});
  const [apiKeys, setApiKeys] = useState<APIKey[]>([]);
  const [usage, setUsage] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [showApiKey, setShowApiKey] = useState<Record<string, boolean>>({});

  useEffect(() => {
    fetchMarketplaceData();
  }, []);

  const fetchMarketplaceData = async () => {
    try {
      const [pricingResponse, keysResponse, usageResponse] = await Promise.all([
        fetch('/api/marketplace/pricing'),
        fetch('/api/marketplace/keys'),
        fetch('/api/marketplace/usage')
      ]);

      const pricingData = await pricingResponse.json();
      const keysData = await keysResponse.json();
      const usageData = await usageResponse.json();

      setEndpoints(pricingData.endpoints || {});
      setApiKeys(keysData.api_keys || []);
      setUsage(usageData);
    } catch (error) {
      console.error('Error fetching marketplace data:', error);
    } finally {
      setLoading(false);
    }
  };

  const createApiKey = async (name: string, description: string, scopes: string[]) => {
    try {
      const response = await fetch('/api/marketplace/keys', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, description, scopes }),
      });

      if (response.ok) {
        fetchMarketplaceData(); // Refresh data
      }
    } catch (error) {
      console.error('Error creating API key:', error);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  const formatPrice = (price: string) => {
    return `$${parseFloat(price).toFixed(3)}`;
  };

  const getCategoryIcon = (category: string) => {
    switch (category) {
      case 'ai_predictions':
      case 'ai_trading':
      case 'ai_analysis':
      case 'ai_portfolio':
      case 'ai_risk':
        return <Brain className="w-5 h-5" />;
      case 'market_data':
        return <BarChart3 className="w-5 h-5" />;
      case 'trading_execution':
      case 'trading_simulation':
        return <TrendingUp className="w-5 h-5" />;
      default:
        return <Code className="w-5 h-5" />;
    }
  };

  const getCategoryColor = (category: string) => {
    switch (category) {
      case 'ai_predictions':
      case 'ai_trading':
      case 'ai_analysis':
      case 'ai_portfolio':
      case 'ai_risk':
        return 'bg-purple-100 text-purple-800 border-purple-200';
      case 'market_data':
        return 'bg-blue-100 text-blue-800 border-blue-200';
      case 'trading_execution':
      case 'trading_simulation':
        return 'bg-green-100 text-green-800 border-green-200';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  if (loading) {
    return (
      <div className="p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-200 rounded w-1/3"></div>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="h-48 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="text-center space-y-4">
        <h1 className="text-4xl font-bold bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent">
          API Marketplace
        </h1>
        <p className="text-xl text-gray-600 max-w-3xl mx-auto">
          Access powerful AI trading algorithms and market data through our API. 
          Pay only for what you use with transparent, per-request pricing.
        </p>
        
        {/* Key Features */}
        <div className="flex justify-center space-x-8 mt-6">
          <div className="flex items-center space-x-2">
            <Brain className="w-5 h-5 text-purple-600" />
            <span className="text-sm font-medium">85%+ AI Accuracy</span>
          </div>
          <div className="flex items-center space-x-2">
            <Zap className="w-5 h-5 text-yellow-600" />
            <span className="text-sm font-medium">Sub-100ms Response</span>
          </div>
          <div className="flex items-center space-x-2">
            <Shield className="w-5 h-5 text-green-600" />
            <span className="text-sm font-medium">Enterprise Security</span>
          </div>
        </div>
      </div>

      <Tabs defaultValue="endpoints" className="space-y-6">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="endpoints">API Endpoints</TabsTrigger>
          <TabsTrigger value="keys">API Keys</TabsTrigger>
          <TabsTrigger value="usage">Usage & Billing</TabsTrigger>
        </TabsList>

        {/* API Endpoints Tab */}
        <TabsContent value="endpoints" className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {Object.entries(endpoints).map(([id, endpoint]) => (
              <Card key={id} className="hover:shadow-lg transition-shadow">
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      {getCategoryIcon(endpoint.category)}
                      <CardTitle className="text-lg">{endpoint.endpoint}</CardTitle>
                    </div>
                    <Badge className={getCategoryColor(endpoint.category)}>
                      {endpoint.category.replace('_', ' ')}
                    </Badge>
                  </div>
                  <CardDescription>{endpoint.description}</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-600">Price per call:</span>
                      <span className="text-lg font-bold text-green-600">
                        {formatPrice(endpoint.price_per_call)}
                      </span>
                    </div>
                    
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-600">Rate limit:</span>
                      <span className="text-sm font-medium">
                        {endpoint.rate_limit}/min
                      </span>
                    </div>
                    
                    {endpoint.requires_plan && (
                      <Badge variant="outline" className="w-full justify-center">
                        Requires Active Subscription
                      </Badge>
                    )}
                    
                    <Button className="w-full" variant="outline">
                      <Code className="w-4 h-4 mr-2" />
                      View Documentation
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          {/* Pricing Tiers */}
          <Card>
            <CardHeader>
              <CardTitle>Volume Discounts</CardTitle>
              <CardDescription>Save more with higher usage volumes</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="text-center p-4 border rounded-lg">
                  <div className="text-2xl font-bold text-blue-600">5%</div>
                  <div className="text-sm text-gray-600">$100+ monthly spend</div>
                </div>
                <div className="text-center p-4 border rounded-lg">
                  <div className="text-2xl font-bold text-green-600">10%</div>
                  <div className="text-sm text-gray-600">$500+ monthly spend</div>
                </div>
                <div className="text-center p-4 border rounded-lg">
                  <div className="text-2xl font-bold text-purple-600">15%</div>
                  <div className="text-sm text-gray-600">$1000+ monthly spend</div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* API Keys Tab */}
        <TabsContent value="keys" className="space-y-6">
          <div className="flex justify-between items-center">
            <h2 className="text-2xl font-bold">API Keys</h2>
            <Button onClick={() => {/* Open create key modal */}}>
              <Key className="w-4 h-4 mr-2" />
              Create New Key
            </Button>
          </div>

          <div className="grid grid-cols-1 gap-4">
            {apiKeys.map((key) => (
              <Card key={key.key_id}>
                <CardContent className="p-6">
                  <div className="flex justify-between items-start">
                    <div className="space-y-2">
                      <h3 className="text-lg font-semibold">{key.name}</h3>
                      <p className="text-gray-600">{key.description}</p>
                      <div className="flex space-x-2">
                        {key.scopes.map((scope) => (
                          <Badge key={scope} variant="outline">
                            {scope}
                          </Badge>
                        ))}
                      </div>
                    </div>
                    
                    <div className="text-right space-y-2">
                      <div className="flex items-center space-x-2">
                        <Input
                          type={showApiKey[key.key_id] ? "text" : "password"}
                          value="aacb_1234567890abcdef1234567890"
                          readOnly
                          className="w-64"
                        />
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => setShowApiKey(prev => ({
                            ...prev,
                            [key.key_id]: !prev[key.key_id]
                          }))}
                        >
                          {showApiKey[key.key_id] ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                        </Button>
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => copyToClipboard("aacb_1234567890abcdef1234567890")}
                        >
                          <Copy className="w-4 h-4" />
                        </Button>
                      </div>
                      
                      <div className="text-sm text-gray-500">
                        <div>Created: {new Date(key.created_at).toLocaleDateString()}</div>
                        <div>Last used: {new Date(key.last_used).toLocaleDateString()}</div>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        {/* Usage & Billing Tab */}
        <TabsContent value="usage" className="space-y-6">
          {usage && (
            <>
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <Card>
                  <CardContent className="p-6">
                    <div className="flex items-center space-x-2">
                      <BarChart3 className="w-5 h-5 text-blue-600" />
                      <div>
                        <div className="text-2xl font-bold">{usage.total_requests?.toLocaleString() || 0}</div>
                        <div className="text-sm text-gray-600">Total Requests</div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
                
                <Card>
                  <CardContent className="p-6">
                    <div className="flex items-center space-x-2">
                      <DollarSign className="w-5 h-5 text-green-600" />
                      <div>
                        <div className="text-2xl font-bold">${usage.total_cost?.toFixed(2) || '0.00'}</div>
                        <div className="text-sm text-gray-600">Total Cost</div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
                
                <Card>
                  <CardContent className="p-6">
                    <div className="flex items-center space-x-2">
                      <TrendingUp className="w-5 h-5 text-purple-600" />
                      <div>
                        <div className="text-2xl font-bold">
                          {Object.keys(usage.endpoint_usage || {}).length}
                        </div>
                        <div className="text-sm text-gray-600">Endpoints Used</div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
                
                <Card>
                  <CardContent className="p-6">
                    <div className="flex items-center space-x-2">
                      <Zap className="w-5 h-5 text-yellow-600" />
                      <div>
                        <div className="text-2xl font-bold">
                          ${(usage.total_cost * 30 / new Date().getDate()).toFixed(0)}
                        </div>
                        <div className="text-sm text-gray-600">Monthly Projection</div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>

              {/* Endpoint Usage Breakdown */}
              <Card>
                <CardHeader>
                  <CardTitle>Endpoint Usage Breakdown</CardTitle>
                  <CardDescription>Detailed usage by endpoint for current billing period</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {Object.entries(usage.endpoint_usage || {}).map(([endpoint, data]: [string, any]) => (
                      <div key={endpoint} className="flex justify-between items-center p-4 border rounded-lg">
                        <div>
                          <div className="font-medium">{endpoint}</div>
                          <div className="text-sm text-gray-600">
                            {data.request_count} requests • {(data.success_rate * 100).toFixed(1)}% success rate
                          </div>
                        </div>
                        <div className="text-right">
                          <div className="font-bold">${data.total_cost.toFixed(3)}</div>
                          <div className="text-sm text-gray-600">
                            {data.avg_duration}ms avg
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default APIMarketplace;
