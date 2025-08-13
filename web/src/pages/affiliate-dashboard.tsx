import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { 
  DollarSign, 
  Users, 
  MousePointer, 
  TrendingUp,
  Copy,
  Share2,
  BarChart3,
  Gift,
  Star,
  Calendar,
  ExternalLink,
  Download
} from 'lucide-react';

interface AffiliateDashboard {
  affiliate_info: {
    affiliate_code: string;
    status: string;
    tier: string;
    commission_rate: string;
    member_since: string;
  };
  performance_summary: {
    total_clicks: number;
    total_signups: number;
    total_conversions: number;
    conversion_rate: string;
    total_commissions: number;
    unpaid_commissions: number;
    this_month_earnings: number;
    last_payout: string;
    next_payout: string;
  };
  recent_referrals: Array<{
    user_id: string;
    conversion_type: string;
    conversion_value: number;
    commission: number;
    status: string;
    date: string;
  }>;
  affiliate_links: {
    main_site: string;
    beta_signup: string;
    api_docs: string;
    pricing: string;
  };
}

const AffiliateDashboard = () => {
  const [dashboard, setDashboard] = useState<AffiliateDashboard | null>(null);
  const [stats, setStats] = useState<any>(null);
  const [referrals, setReferrals] = useState<any[]>([]);
  const [payouts, setPayouts] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [copiedLink, setCopiedLink] = useState<string>('');

  useEffect(() => {
    fetchAffiliateData();
  }, []);

  const fetchAffiliateData = async () => {
    try {
      const [dashboardResponse, statsResponse, referralsResponse, payoutsResponse] = await Promise.all([
        fetch('/api/affiliate/dashboard'),
        fetch('/api/affiliate/stats'),
        fetch('/api/affiliate/referrals'),
        fetch('/api/affiliate/payouts')
      ]);

      const dashboardData = await dashboardResponse.json();
      const statsData = await statsResponse.json();
      const referralsData = await referralsResponse.json();
      const payoutsData = await payoutsResponse.json();

      setDashboard(dashboardData);
      setStats(statsData);
      setReferrals(referralsData.referrals || []);
      setPayouts(payoutsData);
    } catch (error) {
      console.error('Error fetching affiliate data:', error);
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = (text: string, linkType: string) => {
    navigator.clipboard.writeText(text);
    setCopiedLink(linkType);
    setTimeout(() => setCopiedLink(''), 2000);
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
    }).format(amount);
  };

  const getTierColor = (tier: string) => {
    switch (tier.toLowerCase()) {
      case 'bronze': return 'bg-orange-100 text-orange-800 border-orange-200';
      case 'silver': return 'bg-gray-100 text-gray-800 border-gray-200';
      case 'gold': return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      case 'platinum': return 'bg-purple-100 text-purple-800 border-purple-200';
      case 'diamond': return 'bg-blue-100 text-blue-800 border-blue-200';
      default: return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  if (loading) {
    return (
      <div className="p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-200 rounded w-1/3"></div>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="h-32 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Affiliate Dashboard</h1>
          <p className="text-gray-600">Track your referrals and earnings</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline">
            <Download className="w-4 h-4 mr-2" />
            Export Report
          </Button>
          <Button>
            <Share2 className="w-4 h-4 mr-2" />
            Share Links
          </Button>
        </div>
      </div>

      {/* Affiliate Info */}
      {dashboard && (
        <Card>
          <CardHeader>
            <div className="flex justify-between items-center">
              <div>
                <CardTitle className="flex items-center gap-2">
                  Affiliate Code: {dashboard.affiliate_info.affiliate_code}
                  <Badge className={getTierColor(dashboard.affiliate_info.tier)}>
                    {dashboard.affiliate_info.tier}
                  </Badge>
                </CardTitle>
                <CardDescription>
                  Member since {dashboard.affiliate_info.member_since} • {dashboard.affiliate_info.commission_rate} commission rate
                </CardDescription>
              </div>
              <Badge variant={dashboard.affiliate_info.status === 'active' ? 'default' : 'secondary'}>
                {dashboard.affiliate_info.status}
              </Badge>
            </div>
          </CardHeader>
        </Card>
      )}

      {/* Performance Metrics */}
      {dashboard && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Clicks</CardTitle>
              <MousePointer className="h-4 w-4 text-blue-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{dashboard.performance_summary.total_clicks.toLocaleString()}</div>
              <p className="text-xs text-gray-600">
                {dashboard.performance_summary.conversion_rate} conversion rate
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Conversions</CardTitle>
              <Users className="h-4 w-4 text-green-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{dashboard.performance_summary.total_conversions}</div>
              <p className="text-xs text-gray-600">
                {dashboard.performance_summary.total_signups} total signups
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Earnings</CardTitle>
              <DollarSign className="h-4 w-4 text-purple-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-green-600">
                {formatCurrency(dashboard.performance_summary.total_commissions)}
              </div>
              <p className="text-xs text-gray-600">
                All-time commissions
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Unpaid Balance</CardTitle>
              <TrendingUp className="h-4 w-4 text-orange-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-orange-600">
                {formatCurrency(dashboard.performance_summary.unpaid_commissions)}
              </div>
              <p className="text-xs text-gray-600">
                Next payout: {dashboard.performance_summary.next_payout}
              </p>
            </CardContent>
          </Card>
        </div>
      )}

      <Tabs defaultValue="overview" className="space-y-6">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="links">Affiliate Links</TabsTrigger>
          <TabsTrigger value="referrals">Referrals</TabsTrigger>
          <TabsTrigger value="payouts">Payouts</TabsTrigger>
        </TabsList>

        {/* Overview Tab */}
        <TabsContent value="overview" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Recent Referrals */}
            <Card>
              <CardHeader>
                <CardTitle>Recent Referrals</CardTitle>
                <CardDescription>Your latest successful referrals</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {dashboard?.recent_referrals.map((referral, index) => (
                    <div key={index} className="flex justify-between items-center p-3 border rounded-lg">
                      <div>
                        <div className="font-medium">{referral.conversion_type}</div>
                        <div className="text-sm text-gray-600">
                          {referral.date} • {formatCurrency(referral.conversion_value)}
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="font-bold text-green-600">
                          {formatCurrency(referral.commission)}
                        </div>
                        <Badge variant={referral.status === 'confirmed' ? 'default' : 'secondary'}>
                          {referral.status}
                        </Badge>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Performance Stats */}
            <Card>
              <CardHeader>
                <CardTitle>Performance Analytics</CardTitle>
                <CardDescription>Your conversion metrics</CardDescription>
              </CardHeader>
              <CardContent>
                {stats && (
                  <div className="space-y-4">
                    <div className="flex justify-between">
                      <span className="text-sm">Conversion Rate</span>
                      <span className="font-medium">{stats.metrics.conversion_rate}%</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-sm">Average Order Value</span>
                      <span className="font-medium">{formatCurrency(stats.metrics.average_order_value)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-sm">Revenue Trend</span>
                      <span className="font-medium text-green-600">{stats.trends.revenue_trend}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-sm">Click Trend</span>
                      <span className="font-medium text-blue-600">{stats.trends.clicks_trend}</span>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Commission Tiers */}
          <Card>
            <CardHeader>
              <CardTitle>Commission Tiers</CardTitle>
              <CardDescription>Earn more with higher referral volumes</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
                {[
                  { tier: 'Bronze', rate: '15%', min: 0, current: dashboard?.affiliate_info.tier === 'Bronze' },
                  { tier: 'Silver', rate: '20%', min: 10, current: dashboard?.affiliate_info.tier === 'Silver' },
                  { tier: 'Gold', rate: '25%', min: 50, current: dashboard?.affiliate_info.tier === 'Gold' },
                  { tier: 'Platinum', rate: '30%', min: 100, current: dashboard?.affiliate_info.tier === 'Platinum' },
                  { tier: 'Diamond', rate: '35%', min: 500, current: dashboard?.affiliate_info.tier === 'Diamond' },
                ].map((tier) => (
                  <div key={tier.tier} className={`p-4 border rounded-lg text-center ${
                    tier.current ? 'border-blue-500 bg-blue-50' : 'border-gray-200'
                  }`}>
                    <div className="font-bold text-lg">{tier.tier}</div>
                    <div className="text-2xl font-bold text-green-600">{tier.rate}</div>
                    <div className="text-sm text-gray-600">{tier.min}+ referrals</div>
                    {tier.current && (
                      <Badge className="mt-2">Current Tier</Badge>
                    )}
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Affiliate Links Tab */}
        <TabsContent value="links" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Your Affiliate Links</CardTitle>
              <CardDescription>Share these links to earn commissions</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {dashboard && Object.entries(dashboard.affiliate_links).map(([key, url]) => (
                  <div key={key} className="flex items-center space-x-4 p-4 border rounded-lg">
                    <div className="flex-1">
                      <div className="font-medium capitalize">{key.replace('_', ' ')}</div>
                      <div className="text-sm text-gray-600 font-mono">{url}</div>
                    </div>
                    <div className="flex space-x-2">
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => copyToClipboard(url, key)}
                      >
                        <Copy className="w-4 h-4" />
                        {copiedLink === key ? 'Copied!' : 'Copy'}
                      </Button>
                      <Button size="sm" variant="outline" asChild>
                        <a href={url} target="_blank" rel="noopener noreferrer">
                          <ExternalLink className="w-4 h-4" />
                        </a>
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Link Generator */}
          <Card>
            <CardHeader>
              <CardTitle>Custom Link Generator</CardTitle>
              <CardDescription>Create custom tracking links for specific campaigns</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <label className="text-sm font-medium">Page</label>
                  <Input placeholder="pricing" />
                </div>
                <div>
                  <label className="text-sm font-medium">Campaign</label>
                  <Input placeholder="social_media" />
                </div>
                <div>
                  <label className="text-sm font-medium">Source</label>
                  <Input placeholder="twitter" />
                </div>
              </div>
              <Button className="mt-4">
                Generate Custom Link
              </Button>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Referrals Tab */}
        <TabsContent value="referrals" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Referral History</CardTitle>
              <CardDescription>Track all your referrals and their status</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {referrals.map((referral) => (
                  <div key={referral.id} className="grid grid-cols-5 gap-4 p-4 border rounded-lg">
                    <div>
                      <div className="font-medium">{referral.conversion_type}</div>
                      <div className="text-sm text-gray-600">{referral.date}</div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">Value</div>
                      <div className="font-medium">{formatCurrency(referral.conversion_value)}</div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">Commission</div>
                      <div className="font-medium text-green-600">{formatCurrency(referral.commission)}</div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">Source</div>
                      <div className="font-medium">{referral.source}</div>
                    </div>
                    <div>
                      <Badge variant={referral.status === 'confirmed' ? 'default' : 'secondary'}>
                        {referral.status}
                      </Badge>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Payouts Tab */}
        <TabsContent value="payouts" className="space-y-6">
          {payouts && (
            <>
              <Card>
                <CardHeader>
                  <CardTitle>Next Payout</CardTitle>
                  <CardDescription>Your upcoming commission payment</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-3 gap-4">
                    <div>
                      <div className="text-sm text-gray-600">Estimated Amount</div>
                      <div className="text-2xl font-bold text-green-600">
                        {formatCurrency(payouts.next_payout.estimated_amount)}
                      </div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">Payout Date</div>
                      <div className="font-medium">{payouts.next_payout.payout_date}</div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">Status</div>
                      <Badge>{payouts.next_payout.status}</Badge>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Payout History</CardTitle>
                  <CardDescription>Your commission payment history</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {payouts.payouts.map((payout: any) => (
                      <div key={payout.id} className="flex justify-between items-center p-4 border rounded-lg">
                        <div>
                          <div className="font-medium">{payout.period}</div>
                          <div className="text-sm text-gray-600">
                            {payout.payment_method} • {payout.completed_at}
                          </div>
                        </div>
                        <div className="text-right">
                          <div className="font-bold">{formatCurrency(payout.net_amount)}</div>
                          <Badge variant={payout.status === 'completed' ? 'default' : 'secondary'}>
                            {payout.status}
                          </Badge>
                        </div>
                      </div>
                    ))}
                  </div>
                  <div className="mt-4 p-4 bg-gray-50 rounded-lg">
                    <div className="text-sm text-gray-600">Total Paid</div>
                    <div className="text-xl font-bold">{formatCurrency(payouts.total_paid)}</div>
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

export default AffiliateDashboard;
