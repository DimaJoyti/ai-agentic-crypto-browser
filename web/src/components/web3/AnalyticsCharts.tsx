'use client'

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  PieChart,
  Pie,
  Cell,
  LineChart,
  Line,
  Area,
  AreaChart,
  ResponsiveContainer
} from 'recharts'
import { TransactionAnalytics } from '@/lib/transaction-analytics'

const COLORS = [
  '#3b82f6', // blue
  '#10b981', // green
  '#f59e0b', // yellow
  '#ef4444', // red
  '#8b5cf6', // purple
  '#06b6d4', // cyan
  '#84cc16', // lime
  '#f97316', // orange
  '#ec4899', // pink
  '#6b7280'  // gray
]

interface ChartProps {
  data: any[]
  height?: number
}

export const AnalyticsCharts = {
  StatusPieChart: ({ data, height = 300 }: ChartProps) => (
    <ResponsiveContainer width="100%" height={height}>
      <PieChart>
        <Pie
          data={data}
          cx="50%"
          cy="50%"
          labelLine={false}
          label={({ name, percent }) => `${name} ${((percent || 0) * 100).toFixed(0)}%`}
          outerRadius={80}
          fill="#8884d8"
          dataKey="value"
        >
          {data.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={entry.color || COLORS[index % COLORS.length]} />
          ))}
        </Pie>
        <Tooltip />
      </PieChart>
    </ResponsiveContainer>
  ),

  TimeSeriesChart: ({ data, height = 400 }: ChartProps) => (
    <ResponsiveContainer width="100%" height={height}>
      <AreaChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="date" />
        <YAxis yAxisId="left" />
        <YAxis yAxisId="right" orientation="right" />
        <Tooltip 
          formatter={(value: any, name: string) => {
            if (name === 'volume') return [TransactionAnalytics.formatCurrency(value), 'Volume (ETH)']
            if (name === 'gasUsed') return [TransactionAnalytics.formatGas(value), 'Avg Gas']
            if (name === 'successRate') return [`${value.toFixed(1)}%`, 'Success Rate']
            return [value, name]
          }}
        />
        <Legend />
        <Area
          yAxisId="left"
          type="monotone"
          dataKey="transactions"
          stackId="1"
          stroke="#3b82f6"
          fill="#3b82f6"
          fillOpacity={0.6}
          name="Transactions"
        />
        <Line
          yAxisId="right"
          type="monotone"
          dataKey="volume"
          stroke="#10b981"
          strokeWidth={2}
          name="Volume"
        />
      </AreaChart>
    </ResponsiveContainer>
  ),

  ChainBarChart: ({ data, height = 400 }: ChartProps) => (
    <ResponsiveContainer width="100%" height={height}>
      <BarChart data={data} layout="horizontal">
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis type="number" />
        <YAxis dataKey="chainName" type="category" width={100} />
        <Tooltip 
          formatter={(value: any, name: string) => {
            if (name === 'volume') return [TransactionAnalytics.formatCurrency(value), 'Volume (ETH)']
            if (name === 'gasUsed') return [TransactionAnalytics.formatGas(value), 'Avg Gas']
            if (name === 'successRate') return [`${value.toFixed(1)}%`, 'Success Rate']
            return [value, name]
          }}
        />
        <Legend />
        <Bar dataKey="transactionCount" fill="#3b82f6" name="Transactions" />
        <Bar dataKey="volume" fill="#10b981" name="Volume" />
      </BarChart>
    </ResponsiveContainer>
  ),

  ChainPieChart: ({ data, height = 300 }: ChartProps) => (
    <ResponsiveContainer width="100%" height={height}>
      <PieChart>
        <Pie
          data={data}
          cx="50%"
          cy="50%"
          labelLine={false}
          label={({ chainName, percent }) => `${chainName} ${((percent || 0) * 100).toFixed(0)}%`}
          outerRadius={80}
          fill="#8884d8"
          dataKey="transactionCount"
        >
          {data.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
          ))}
        </Pie>
        <Tooltip formatter={(value) => [value, 'Transactions']} />
      </PieChart>
    </ResponsiveContainer>
  ),

  TypeBarChart: ({ data, height = 400 }: ChartProps) => (
    <ResponsiveContainer width="100%" height={height}>
      <BarChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis
          dataKey="type"
          angle={-45}
          textAnchor="end"
          height={80}
          tickFormatter={(value) => value.replace('_', ' ')}
        />
        <YAxis />
        <Tooltip 
          formatter={(value: any, name: string) => {
            if (name === 'volume') return [TransactionAnalytics.formatCurrency(value), 'Volume (ETH)']
            if (name === 'averageGasUsed') return [TransactionAnalytics.formatGas(value), 'Avg Gas']
            if (name === 'successRate') return [`${value.toFixed(1)}%`, 'Success Rate']
            return [value, name]
          }}
          labelFormatter={(label) => label.replace('_', ' ')}
        />
        <Legend />
        <Bar dataKey="count" fill="#3b82f6" name="Count" />
        <Bar dataKey="volume" fill="#10b981" name="Volume" />
      </BarChart>
    </ResponsiveContainer>
  ),

  TypePieChart: ({ data, height = 300 }: ChartProps) => (
    <ResponsiveContainer width="100%" height={height}>
      <PieChart>
        <Pie
          data={data}
          cx="50%"
          cy="50%"
          labelLine={false}
          label={({ type, percent }) => `${type.replace('_', ' ')} ${((percent || 0) * 100).toFixed(0)}%`}
          outerRadius={80}
          fill="#8884d8"
          dataKey="count"
        >
          {data.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
          ))}
        </Pie>
        <Tooltip 
          formatter={(value) => [value, 'Transactions']}
          labelFormatter={(label) => label.replace('_', ' ')}
        />
      </PieChart>
    </ResponsiveContainer>
  ),

  SuccessRateTrend: ({ data, height = 300 }: ChartProps) => (
    <ResponsiveContainer width="100%" height={height}>
      <LineChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="date" />
        <YAxis domain={[0, 100]} />
        <Tooltip formatter={(value: any) => [`${value.toFixed(1)}%`, 'Success Rate']} />
        <Line
          type="monotone"
          dataKey="successRate"
          stroke="#10b981"
          strokeWidth={2}
          dot={{ fill: '#10b981' }}
        />
      </LineChart>
    </ResponsiveContainer>
  ),

  GasUsageTrend: ({ data, height = 300 }: ChartProps) => (
    <ResponsiveContainer width="100%" height={height}>
      <LineChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="date" />
        <YAxis />
        <Tooltip formatter={(value: any) => [TransactionAnalytics.formatGas(value), 'Avg Gas Used']} />
        <Line
          type="monotone"
          dataKey="gasUsed"
          stroke="#f59e0b"
          strokeWidth={2}
          dot={{ fill: '#f59e0b' }}
        />
      </LineChart>
    </ResponsiveContainer>
  ),

  VolumeTrend: ({ data, height = 400 }: ChartProps) => (
    <ResponsiveContainer width="100%" height={height}>
      <AreaChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="date" />
        <YAxis />
        <Tooltip formatter={(value: any) => [TransactionAnalytics.formatCurrency(value), 'Volume (ETH)']} />
        <Area
          type="monotone"
          dataKey="volume"
          stroke="#3b82f6"
          fill="#3b82f6"
          fillOpacity={0.6}
        />
      </AreaChart>
    </ResponsiveContainer>
  ),

  TransactionHeatmap: ({ data, height = 400 }: ChartProps) => (
    <ResponsiveContainer width="100%" height={height}>
      <BarChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="hour" />
        <YAxis />
        <Tooltip />
        <Bar dataKey="transactions" fill="#3b82f6" />
      </BarChart>
    </ResponsiveContainer>
  ),

  GasEfficiencyChart: ({ data, height = 400 }: ChartProps) => (
    <ResponsiveContainer width="100%" height={height}>
      <BarChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="chainName" />
        <YAxis />
        <Tooltip 
          formatter={(value: any, name: string) => {
            if (name === 'gasUsed') return [TransactionAnalytics.formatGas(value), 'Avg Gas Used']
            if (name === 'averageConfirmationTime') return [TransactionAnalytics.formatTime(value), 'Avg Time']
            return [value, name]
          }}
        />
        <Legend />
        <Bar dataKey="gasUsed" fill="#f59e0b" name="Gas Used" />
        <Bar dataKey="averageConfirmationTime" fill="#8b5cf6" name="Confirmation Time" />
      </BarChart>
    </ResponsiveContainer>
  )
}
