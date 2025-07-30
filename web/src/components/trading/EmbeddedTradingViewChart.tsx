import React from 'react';

interface EmbeddedTradingViewChartProps {
  symbol: string;
  theme?: 'light' | 'dark';
  interval?: string;
  height?: number;
  className?: string;
}

export const EmbeddedTradingViewChart: React.FC<EmbeddedTradingViewChartProps> = ({
  symbol = 'BTCUSDT',
  theme = 'dark',
  interval = '1D',
  height = 500,
  className = ''
}) => {
  // Create the TradingView URL with parameters
  const tradingViewUrl = `https://www.tradingview.com/widgetembed/?frameElementId=tradingview_chart&symbol=BINANCE%3A${symbol}&interval=${interval}&hidesidetoolbar=1&hidetoptoolbar=0&symboledit=1&saveimage=1&toolbarbg=${theme === 'dark' ? '000000' : 'f1f3f6'}&studies=%5B%5D&theme=${theme}&style=1&timezone=Etc%2FUTC&studies_overrides=%7B%7D&overrides=%7B%7D&enabled_features=%5B%5D&disabled_features=%5B%5D&locale=en&utm_source=localhost&utm_medium=widget_new&utm_campaign=chart&utm_term=BINANCE%3A${symbol}`;

  return (
    <div className={`tradingview-widget-container ${className}`}>
      <div className="tradingview-widget-container__widget">
        <iframe
          id="tradingview_chart"
          src={tradingViewUrl}
          style={{
            width: '100%',
            height: `${height}px`,
            border: 'none',
            borderRadius: '8px',
            backgroundColor: theme === 'dark' ? '#1a1a1a' : '#ffffff'
          }}
          frameBorder="0"
          allowTransparency={true}
          scrolling="no"
          allowFullScreen={true}
          title={`${symbol} TradingView Chart`}
        />
      </div>
      <div className="tradingview-widget-copyright">
        <a 
          href={`https://www.tradingview.com/symbols/${symbol}/?exchange=BINANCE`}
          rel="noopener nofollow" 
          target="_blank"
          className="text-xs text-gray-500 hover:text-gray-400"
        >
          Track all markets on TradingView
        </a>
      </div>
    </div>
  );
};

export default EmbeddedTradingViewChart;
