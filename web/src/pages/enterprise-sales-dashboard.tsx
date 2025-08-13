import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Progress } from '@/components/ui/progress';
import { 
  DollarSign, 
  Users, 
  Target, 
  TrendingUp,
  Calendar,
  Phone,
  Mail,
  Building,
  Award,
  BarChart3,
  PieChart,
  Activity,
  Clock,
  CheckCircle,
  AlertCircle
} from 'lucide-react';

interface SalesDashboard {
  overview: {
    total_leads: number;
    qualified_leads: number;
    active_deals: number;
    pipeline_value: number;
    closed_this_month: number;
    quota_achievement: string;
  };
  pipeline_by_stage: Array<{
    stage: string;
    count: number;
    value: number;
  }>;
  top_opportunities: Array<{
    company: string;
    value: number;
    stage: string;
    probability: number;
    close_date: string;
  }>;
  recent_activities: Array<{
    type: string;
    company: string;
    date: string;
    outcome: string;
    next_steps: string;
  }>;
  performance_metrics: {
    conversion_rate: string;
    avg_deal_size: string;
    avg_sales_cycle: string;
    activities_this_week: number;
  };
}

const EnterpriseSalesDashboard = () => {
  const [dashboard, setDashboard] = useState<SalesDashboard | null>(null);
  const [leads, setLeads] = useState<any[]>([]);
  const [deals, setDeals] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchSalesData();
  }, []);

  const fetchSalesData = async () => {
    try {
      const [dashboardResponse, leadsResponse, dealsResponse] = await Promise.all([
        fetch('/api/sales/dashboard'),
        fetch('/api/sales/leads?limit=10'),
        fetch('/api/sales/deals?limit=10')
      ]);

      const dashboardData = await dashboardResponse.json();
      const leadsData = await leadsResponse.json();
      const dealsData = await dealsResponse.json();

      setDashboard(dashboardData);
      setLeads(leadsData.leads || []);
      setDeals(dealsData.deals || []);
    } catch (error) {
      console.error('Error fetching sales data:', error);
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  const getStageColor = (stage: string) => {
    switch (stage.toLowerCase()) {
      case 'discovery': return 'bg-blue-100 text-blue-800 border-blue-200';
      case 'demo': return 'bg-purple-100 text-purple-800 border-purple-200';
      case 'proposal': return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      case 'negotiation': return 'bg-orange-100 text-orange-800 border-orange-200';
      case 'contract': return 'bg-green-100 text-green-800 border-green-200';
      default: return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  const getPriorityColor = (priority: string) => {
    switch (priority.toLowerCase()) {
      case 'critical': return 'bg-red-100 text-red-800 border-red-200';
      case 'high': return 'bg-orange-100 text-orange-800 border-orange-200';
      case 'medium': return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      case 'low': return 'bg-green-100 text-green-800 border-green-200';
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
          <h1 className="text-3xl font-bold">Enterprise Sales Dashboard</h1>
          <p className="text-gray-600">Manage your enterprise sales pipeline</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline">
            <BarChart3 className="w-4 h-4 mr-2" />
            Reports
          </Button>
          <Button>
            <Users className="w-4 h-4 mr-2" />
            New Lead
          </Button>
        </div>
      </div>

      {/* Key Metrics */}
      {dashboard && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Pipeline Value</CardTitle>
              <DollarSign className="h-4 w-4 text-green-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-green-600">
                {formatCurrency(dashboard.overview.pipeline_value)}
              </div>
              <p className="text-xs text-gray-600">
                {dashboard.overview.active_deals} active deals
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Qualified Leads</CardTitle>
              <Users className="h-4 w-4 text-blue-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-blue-600">
                {dashboard.overview.qualified_leads}
              </div>
              <p className="text-xs text-gray-600">
                of {dashboard.overview.total_leads} total leads
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Monthly Closed</CardTitle>
              <Target className="h-4 w-4 text-purple-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-purple-600">
                {formatCurrency(dashboard.overview.closed_this_month)}
              </div>
              <p className="text-xs text-gray-600">
                This month's revenue
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Quota Achievement</CardTitle>
              <Award className="h-4 w-4 text-orange-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-orange-600">
                {dashboard.overview.quota_achievement}
              </div>
              <Progress value={parseInt(dashboard.overview.quota_achievement)} className="mt-2" />
            </CardContent>
          </Card>
        </div>
      )}

      <Tabs defaultValue="pipeline" className="space-y-6">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="pipeline">Pipeline</TabsTrigger>
          <TabsTrigger value="leads">Leads</TabsTrigger>
          <TabsTrigger value="deals">Deals</TabsTrigger>
          <TabsTrigger value="activities">Activities</TabsTrigger>
        </TabsList>

        {/* Pipeline Tab */}
        <TabsContent value="pipeline" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Pipeline by Stage */}
            <Card>
              <CardHeader>
                <CardTitle>Pipeline by Stage</CardTitle>
                <CardDescription>Deals distribution across sales stages</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {dashboard?.pipeline_by_stage.map((stage) => (
                    <div key={stage.stage} className="flex justify-between items-center">
                      <div className="flex items-center space-x-3">
                        <Badge className={getStageColor(stage.stage)}>
                          {stage.stage}
                        </Badge>
                        <span className="text-sm text-gray-600">{stage.count} deals</span>
                      </div>
                      <div className="font-medium">{formatCurrency(stage.value)}</div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Top Opportunities */}
            <Card>
              <CardHeader>
                <CardTitle>Top Opportunities</CardTitle>
                <CardDescription>Highest value deals in pipeline</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {dashboard?.top_opportunities.map((opp, index) => (
                    <div key={index} className="flex justify-between items-center p-3 border rounded-lg">
                      <div>
                        <div className="font-medium">{opp.company}</div>
                        <div className="text-sm text-gray-600">
                          {opp.stage} • {opp.probability}% probability
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="font-bold text-green-600">
                          {formatCurrency(opp.value)}
                        </div>
                        <div className="text-sm text-gray-600">{opp.close_date}</div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Performance Metrics */}
          <Card>
            <CardHeader>
              <CardTitle>Performance Metrics</CardTitle>
              <CardDescription>Key sales performance indicators</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
                <div className="text-center">
                  <div className="text-2xl font-bold text-blue-600">
                    {dashboard?.performance_metrics.conversion_rate}
                  </div>
                  <div className="text-sm text-gray-600">Conversion Rate</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-green-600">
                    {dashboard?.performance_metrics.avg_deal_size}
                  </div>
                  <div className="text-sm text-gray-600">Avg Deal Size</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-purple-600">
                    {dashboard?.performance_metrics.avg_sales_cycle}
                  </div>
                  <div className="text-sm text-gray-600">Avg Sales Cycle</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-orange-600">
                    {dashboard?.performance_metrics.activities_this_week}
                  </div>
                  <div className="text-sm text-gray-600">Activities This Week</div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Leads Tab */}
        <TabsContent value="leads" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Recent Leads</CardTitle>
              <CardDescription>Latest enterprise prospects</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {/* Mock lead data */}
                {[
                  {
                    company: "Crypto Capital Partners",
                    contact: "John Smith",
                    email: "john@cryptocapital.com",
                    type: "Hedge Fund",
                    aum: "$250M",
                    status: "Qualified",
                    priority: "High",
                    next_follow_up: "2024-01-30",
                  },
                  {
                    company: "Digital Asset Fund",
                    contact: "Sarah Johnson",
                    email: "sarah@digitalasset.com",
                    type: "Family Office",
                    aum: "$150M",
                    status: "New",
                    priority: "Medium",
                    next_follow_up: "2024-01-29",
                  },
                ].map((lead, index) => (
                  <div key={index} className="grid grid-cols-6 gap-4 p-4 border rounded-lg">
                    <div>
                      <div className="font-medium">{lead.company}</div>
                      <div className="text-sm text-gray-600">{lead.contact}</div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">Type</div>
                      <div className="font-medium">{lead.type}</div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">AUM</div>
                      <div className="font-medium">{lead.aum}</div>
                    </div>
                    <div>
                      <Badge className={getStageColor(lead.status)}>
                        {lead.status}
                      </Badge>
                    </div>
                    <div>
                      <Badge className={getPriorityColor(lead.priority)}>
                        {lead.priority}
                      </Badge>
                    </div>
                    <div className="flex space-x-2">
                      <Button size="sm" variant="outline">
                        <Phone className="w-4 h-4" />
                      </Button>
                      <Button size="sm" variant="outline">
                        <Mail className="w-4 h-4" />
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Deals Tab */}
        <TabsContent value="deals" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Active Deals</CardTitle>
              <CardDescription>Deals in progress</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {/* Mock deal data */}
                {[
                  {
                    company: "Blockchain Ventures",
                    value: "$180,000",
                    stage: "Proposal",
                    probability: 60,
                    close_date: "2024-02-15",
                    sales_rep: "Alice Johnson",
                  },
                  {
                    company: "Crypto Innovations",
                    value: "$320,000",
                    stage: "Negotiation",
                    probability: 75,
                    close_date: "2024-02-28",
                    sales_rep: "Bob Wilson",
                  },
                ].map((deal, index) => (
                  <div key={index} className="grid grid-cols-6 gap-4 p-4 border rounded-lg">
                    <div>
                      <div className="font-medium">{deal.company}</div>
                      <div className="text-sm text-gray-600">{deal.sales_rep}</div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">Value</div>
                      <div className="font-medium text-green-600">{deal.value}</div>
                    </div>
                    <div>
                      <Badge className={getStageColor(deal.stage)}>
                        {deal.stage}
                      </Badge>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">Probability</div>
                      <div className="font-medium">{deal.probability}%</div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">Close Date</div>
                      <div className="font-medium">{deal.close_date}</div>
                    </div>
                    <div className="flex space-x-2">
                      <Button size="sm" variant="outline">
                        Edit
                      </Button>
                      <Button size="sm">
                        Update
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Activities Tab */}
        <TabsContent value="activities" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Recent Activities</CardTitle>
              <CardDescription>Latest sales activities and touchpoints</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {dashboard?.recent_activities.map((activity, index) => (
                  <div key={index} className="flex items-center space-x-4 p-4 border rounded-lg">
                    <div className="p-2 rounded-full bg-blue-100">
                      {activity.type === 'Demo' && <Activity className="w-4 h-4 text-blue-600" />}
                      {activity.type === 'Call' && <Phone className="w-4 h-4 text-blue-600" />}
                      {activity.type === 'Email' && <Mail className="w-4 h-4 text-blue-600" />}
                      {activity.type === 'Meeting' && <Calendar className="w-4 h-4 text-blue-600" />}
                    </div>
                    <div className="flex-1">
                      <div className="font-medium">{activity.type} with {activity.company}</div>
                      <div className="text-sm text-gray-600">
                        {activity.date} • Outcome: {activity.outcome}
                      </div>
                      <div className="text-sm text-gray-600">
                        Next: {activity.next_steps}
                      </div>
                    </div>
                    <div className="flex items-center space-x-2">
                      {activity.outcome === 'Positive' && (
                        <CheckCircle className="w-5 h-5 text-green-600" />
                      )}
                      {activity.outcome === 'Negotiation' && (
                        <AlertCircle className="w-5 h-5 text-orange-600" />
                      )}
                      <Clock className="w-4 h-4 text-gray-400" />
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default EnterpriseSalesDashboard;
