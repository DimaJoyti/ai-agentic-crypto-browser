import React, { useEffect, useRef } from 'react';

interface SimpleTradingViewChartProps {
  symbol: string;
  theme?: 'light' | 'dark';
  interval?: string;
  height?: number;
  className?: string;
}

export const SimpleTradingViewChart: React.FC<SimpleTradingViewChartProps> = ({
  symbol = 'BTCUSDT',
  theme = 'dark',
  interval = '1D',
  height = 500,
  className = ''
}) => {
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!containerRef.current) return;

    // Clear any existing content
    containerRef.current.innerHTML = '';

    // Create the script element
    const script = document.createElement('script');
    script.type = 'text/javascript';
    script.src = 'https://s3.tradingview.com/external-embedding/embed-widget-symbol-overview.js';
    script.async = true;

    // Widget configuration
    const config = {
      symbols: [
        [`BINANCE:${symbol}|1D`]
      ],
      chartOnly: false,
      width: '100%',
      height: height,
      locale: 'en',
      colorTheme: theme,
      autosize: true,
      showVolume: false,
      showMA: false,
      hideDateRanges: false,
      hideMarketStatus: false,
      hideSymbolLogo: false,
      scalePosition: 'right',
      scaleMode: 'Normal',
      fontFamily: '-apple-system, BlinkMacSystemFont, Trebuchet MS, Roboto, Ubuntu, sans-serif',
      fontSize: '10',
      noTimeScale: false,
      valuesTracking: '1',
      changeMode: 'price-and-percent',
      chartType: 'area',
      maLineColor: '#2962FF',
      maLineWidth: 1,
      maLength: 9,
      backgroundColor: theme === 'dark' ? 'rgba(26, 26, 26, 1)' : 'rgba(255, 255, 255, 1)',
      lineWidth: 2,
      lineType: 0,
      dateRanges: [
        '1d|1',
        '1m|30',
        '3m|60',
        '12m|1D',
        '60m|1W',
        'all|1M'
      ]
    };

    script.innerHTML = JSON.stringify(config);

    // Append script to container
    containerRef.current.appendChild(script);

    // Cleanup function
    return () => {
      if (containerRef.current) {
        containerRef.current.innerHTML = '';
      }
    };
  }, [symbol, theme, interval, height]);

  return (
    <div className={`tradingview-widget-container ${className}`}>
      <div 
        ref={containerRef}
        className="tradingview-widget"
        style={{ height: `${height}px` }}
      />
      <div className="tradingview-widget-copyright">
        <a 
          href={`https://www.tradingview.com/symbols/${symbol}/?exchange=BINANCE`}
          rel="noopener nofollow" 
          target="_blank"
          className="text-xs text-gray-500 hover:text-gray-400"
        >
          <span className="blue-text">{symbol}</span> by TradingView
        </a>
      </div>
    </div>
  );
};

export default SimpleTradingViewChart;
