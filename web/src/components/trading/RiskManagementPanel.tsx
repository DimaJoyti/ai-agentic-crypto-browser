import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { 
  Shield, 
  AlertTriangle, 
  CheckCircle, 
  TrendingDown,
  DollarSign,
  Activity
} from 'lucide-react';

export const RiskManagementPanel: React.FC = () => {
  const riskLimits = [
    { name: 'Daily Loss Limit', current: 2500, limit: 5000, status: 'OK' },
    { name: 'Position Size Limit', current: 0.8, limit: 1.0, status: 'OK' },
    { name: 'Max Drawdown', current: 1200, limit: 2000, status: 'WARNING' },
    { name: 'Order Rate Limit', current: 45, limit: 60, status: 'OK' }
  ];

  const violations = [
    {
      id: '1',
      type: 'Position Size',
      symbol: 'BTCUSDT',
      severity: 'HIGH',
      message: 'Position size exceeded 80% of limit',
      timestamp: '2024-01-15 14:30:25'
    }
  ];

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'OK': return 'text-green-600';
      case 'WARNING': return 'text-yellow-600';
      case 'CRITICAL': return 'text-red-600';
      default: return 'text-gray-600';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'OK': return <CheckCircle className="w-4 h-4 text-green-600" />;
      case 'WARNING': return <AlertTriangle className="w-4 h-4 text-yellow-600" />;
      case 'CRITICAL': return <AlertTriangle className="w-4 h-4 text-red-600" />;
      default: return <Activity className="w-4 h-4 text-gray-600" />;
    }
  };

  return (
    <div className="space-y-6">
      {/* Risk Status Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Risk Status</CardTitle>
            <Shield className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">HEALTHY</div>
            <p className="text-xs text-muted-foreground">
              All systems operational
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active Violations</CardTitle>
            <AlertTriangle className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-yellow-600">1</div>
            <p className="text-xs text-muted-foreground">
              1 high severity
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Daily Loss</CardTitle>
            <TrendingDown className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">$2,500</div>
            <p className="text-xs text-muted-foreground">
              50% of limit used
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Emergency Stop</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">READY</div>
            <Button variant="destructive" size="sm" className="mt-2">
              EMERGENCY STOP
            </Button>
          </CardContent>
        </Card>
      </div>

      {/* Active Violations */}
      {violations.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <AlertTriangle className="w-5 h-5 mr-2 text-yellow-600" />
              Active Risk Violations
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {violations.map((violation) => (
                <Alert key={violation.id}>
                  <AlertTriangle className="h-4 w-4" />
                  <AlertDescription>
                    <div className="flex items-center justify-between">
                      <div>
                        <strong>{violation.type}</strong> - {violation.symbol}
                        <p className="text-sm mt-1">{violation.message}</p>
                        <p className="text-xs text-muted-foreground mt-1">{violation.timestamp}</p>
                      </div>
                      <Badge variant="destructive">{violation.severity}</Badge>
                    </div>
                  </AlertDescription>
                </Alert>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Risk Limits */}
      <Card>
        <CardHeader>
          <CardTitle>Risk Limits</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-6">
            {riskLimits.map((limit, index) => {
              const percentage = (limit.current / limit.limit) * 100;
              return (
                <div key={index} className="space-y-2">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      {getStatusIcon(limit.status)}
                      <span className="font-medium">{limit.name}</span>
                    </div>
                    <div className="text-right">
                      <span className={`font-semibold ${getStatusColor(limit.status)}`}>
                        {limit.current.toLocaleString()} / {limit.limit.toLocaleString()}
                      </span>
                      <Badge variant="outline" className="ml-2">
                        {limit.status}
                      </Badge>
                    </div>
                  </div>
                  <Progress 
                    value={percentage} 
                    className={`h-2 ${
                      percentage > 80 ? 'bg-red-100' : 
                      percentage > 60 ? 'bg-yellow-100' : 'bg-green-100'
                    }`}
                  />
                  <div className="flex justify-between text-xs text-muted-foreground">
                    <span>0</span>
                    <span>{percentage.toFixed(1)}% used</span>
                    <span>{limit.limit.toLocaleString()}</span>
                  </div>
                </div>
              );
            })}
          </div>
        </CardContent>
      </Card>

      {/* Risk Controls */}
      <Card>
        <CardHeader>
          <CardTitle>Risk Controls</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-3">
              <h4 className="font-semibold">Trading Controls</h4>
              <div className="space-y-2">
                <Button variant="outline" className="w-full justify-start">
                  Halt All Trading
                </Button>
                <Button variant="outline" className="w-full justify-start">
                  Cancel All Orders
                </Button>
                <Button variant="outline" className="w-full justify-start">
                  Close All Positions
                </Button>
              </div>
            </div>
            
            <div className="space-y-3">
              <h4 className="font-semibold">Risk Settings</h4>
              <div className="space-y-2">
                <Button variant="outline" className="w-full justify-start">
                  Adjust Position Limits
                </Button>
                <Button variant="outline" className="w-full justify-start">
                  Modify Risk Parameters
                </Button>
                <Button variant="outline" className="w-full justify-start">
                  Update Stop Loss Rules
                </Button>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};
