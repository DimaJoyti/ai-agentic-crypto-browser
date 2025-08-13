import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { 
  TrendingUp, 
  DollarSign, 
  Users, 
  Target, 
  Calendar,
  ArrowUpRight,
  ArrowDownRight
} from 'lucide-react';

interface RevenueMetrics {
  totalRevenue: number;
  monthlyRevenue: number;
  subscriptionRevenue: number;
  performanceFees: number;
  apiRevenue: number;
  growthRate: number;
  churnRate: number;
  totalUsers: number;
  activeSubscribers: number;
  trialUsers: number;
  enterpriseClients: number;
  conversionRate: number;
  ltv: number;
}

const RevenueDashboard = () => {
  const [metrics, setMetrics] = useState<RevenueMetrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [timeframe, setTimeframe] = useState('30d');

  useEffect(() => {
    fetchMetrics();
  }, [timeframe]);

  const fetchMetrics = async () => {
    try {
      const [revenueResponse, userResponse] = await Promise.all([
        fetch('/api/billing/analytics/revenue'),
        fetch('/api/billing/analytics/users')
      ]);

      const revenueData = await revenueResponse.json();
      const userData = await userResponse.json();

      setMetrics({
        ...revenueData,
        ...userData
      });
    } catch (error) {
      console.error('Error fetching metrics:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-200 rounded w-1/4"></div>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="h-32 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (!metrics) {
    return <div className="p-6">Error loading metrics</div>;
  }

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  const formatPercentage = (rate: number) => {
    return `${(rate * 100).toFixed(1)}%`;
  };

  const revenueCards = [
    {
      title: 'Total Revenue',
      value: formatCurrency(metrics.totalRevenue),
      change: `+${formatPercentage(metrics.growthRate)}`,
      icon: DollarSign,
      color: 'text-green-600',
      bgColor: 'bg-green-50',
      trend: 'up'
    },
    {
      title: 'Monthly Revenue',
      value: formatCurrency(metrics.monthlyRevenue),
      change: '+15.2%',
      icon: TrendingUp,
      color: 'text-blue-600',
      bgColor: 'bg-blue-50',
      trend: 'up'
    },
    {
      title: 'Active Subscribers',
      value: metrics.activeSubscribers.toLocaleString(),
      change: '+8.5%',
      icon: Users,
      color: 'text-purple-600',
      bgColor: 'bg-purple-50',
      trend: 'up'
    },
    {
      title: 'Conversion Rate',
      value: formatPercentage(metrics.conversionRate),
      change: '+2.1%',
      icon: Target,
      color: 'text-orange-600',
      bgColor: 'bg-orange-50',
      trend: 'up'
    }
  ];

  const revenueBreakdown = [
    {
      source: 'Subscriptions',
      amount: metrics.subscriptionRevenue,
      percentage: (metrics.subscriptionRevenue / metrics.totalRevenue) * 100,
      color: 'bg-blue-500'
    },
    {
      source: 'Performance Fees',
      amount: metrics.performanceFees,
      percentage: (metrics.performanceFees / metrics.totalRevenue) * 100,
      color: 'bg-green-500'
    },
    {
      source: 'API Revenue',
      amount: metrics.apiRevenue,
      percentage: (metrics.apiRevenue / metrics.totalRevenue) * 100,
      color: 'bg-purple-500'
    }
  ];

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Revenue Dashboard</h1>
          <p className="text-gray-600">Track your platform's financial performance</p>
        </div>
        <div className="flex gap-2">
          {['7d', '30d', '90d', '1y'].map((period) => (
            <Button
              key={period}
              variant={timeframe === period ? 'default' : 'outline'}
              size="sm"
              onClick={() => setTimeframe(period)}
            >
              {period}
            </Button>
          ))}
        </div>
      </div>

      {/* Key Metrics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {revenueCards.map((card, index) => (
          <Card key={index}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-gray-600">
                {card.title}
              </CardTitle>
              <div className={`p-2 rounded-lg ${card.bgColor}`}>
                <card.icon className={`h-4 w-4 ${card.color}`} />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{card.value}</div>
              <div className="flex items-center text-sm">
                {card.trend === 'up' ? (
                  <ArrowUpRight className="h-4 w-4 text-green-500 mr-1" />
                ) : (
                  <ArrowDownRight className="h-4 w-4 text-red-500 mr-1" />
                )}
                <span className={card.trend === 'up' ? 'text-green-600' : 'text-red-600'}>
                  {card.change}
                </span>
                <span className="text-gray-500 ml-1">vs last month</span>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Revenue Breakdown */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Revenue Breakdown</CardTitle>
            <CardDescription>Revenue by source for the selected period</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {revenueBreakdown.map((item, index) => (
                <div key={index} className="space-y-2">
                  <div className="flex justify-between items-center">
                    <span className="text-sm font-medium">{item.source}</span>
                    <span className="text-sm text-gray-600">
                      {formatCurrency(item.amount)} ({item.percentage.toFixed(1)}%)
                    </span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className={`h-2 rounded-full ${item.color}`}
                      style={{ width: `${item.percentage}%` }}
                    ></div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>User Metrics</CardTitle>
            <CardDescription>User acquisition and retention metrics</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex justify-between items-center">
                <span className="text-sm font-medium">Total Users</span>
                <span className="text-lg font-bold">{metrics.totalUsers.toLocaleString()}</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm font-medium">Active Subscribers</span>
                <span className="text-lg font-bold text-green-600">
                  {metrics.activeSubscribers.toLocaleString()}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm font-medium">Trial Users</span>
                <span className="text-lg font-bold text-blue-600">
                  {metrics.trialUsers.toLocaleString()}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm font-medium">Enterprise Clients</span>
                <span className="text-lg font-bold text-purple-600">
                  {metrics.enterpriseClients.toLocaleString()}
                </span>
              </div>
              <div className="pt-4 border-t">
                <div className="flex justify-between items-center">
                  <span className="text-sm font-medium">Customer LTV</span>
                  <span className="text-lg font-bold">{formatCurrency(metrics.ltv)}</span>
                </div>
                <div className="flex justify-between items-center mt-2">
                  <span className="text-sm font-medium">Churn Rate</span>
                  <span className="text-lg font-bold text-red-600">
                    {formatPercentage(metrics.churnRate)}
                  </span>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Revenue Projections */}
      <Card>
        <CardHeader>
          <CardTitle>Revenue Projections</CardTitle>
          <CardDescription>Projected revenue based on current growth trends</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="text-center">
              <div className="text-2xl font-bold text-blue-600">
                {formatCurrency(metrics.monthlyRevenue * 3)}
              </div>
              <div className="text-sm text-gray-600">3-Month Projection</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-green-600">
                {formatCurrency(metrics.monthlyRevenue * 6)}
              </div>
              <div className="text-sm text-gray-600">6-Month Projection</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-purple-600">
                {formatCurrency(metrics.monthlyRevenue * 12)}
              </div>
              <div className="text-sm text-gray-600">Annual Projection</div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Action Items */}
      <Card>
        <CardHeader>
          <CardTitle>Revenue Optimization</CardTitle>
          <CardDescription>Recommended actions to increase revenue</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <h4 className="font-semibold text-green-600">Opportunities</h4>
              <ul className="text-sm space-y-1">
                <li>• Increase Professional tier conversion (+$50K/month potential)</li>
                <li>• Launch enterprise white-label program (+$100K/month)</li>
                <li>• Expand API marketplace (+$25K/month)</li>
                <li>• Add performance fee tiers (+$75K/month)</li>
              </ul>
            </div>
            <div className="space-y-2">
              <h4 className="font-semibold text-red-600">Risks</h4>
              <ul className="text-sm space-y-1">
                <li>• Churn rate above 5% target</li>
                <li>• Trial to paid conversion below 25%</li>
                <li>• Enterprise pipeline needs attention</li>
                <li>• API usage growth slowing</li>
              </ul>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default RevenueDashboard;
