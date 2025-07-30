import React, { useEffect, useRef, memo } from 'react';

interface TradingViewWidgetProps {
  symbol: string;
  theme?: 'light' | 'dark';
  interval?: string;
  width?: string | number;
  height?: string | number;
  autosize?: boolean;
  className?: string;
}

declare global {
  interface Window {
    TradingView: any;
  }
}

const TradingViewWidget: React.FC<TradingViewWidgetProps> = memo(({
  symbol = 'BINANCE:BTCUSDT',
  theme = 'dark',
  interval = '1D',
  width = '100%',
  height = 500,
  autosize = true,
  className = ''
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const scriptRef = useRef<HTMLScriptElement | null>(null);

  useEffect(() => {
    // Create script element
    const script = document.createElement('script');
    script.src = 'https://s3.tradingview.com/external-embedding/embed-widget-advanced-chart.js';
    script.type = 'text/javascript';
    script.async = true;
    
    // Widget configuration
    script.innerHTML = JSON.stringify({
      autosize: autosize,
      width: autosize ? undefined : width,
      height: autosize ? undefined : height,
      symbol: symbol,
      interval: interval,
      timezone: 'Etc/UTC',
      theme: theme,
      style: '1',
      locale: 'en',
      enable_publishing: false,
      allow_symbol_change: true,
      calendar: false,
      support_host: 'https://www.tradingview.com',
      container_id: containerRef.current?.id || 'tradingview-widget',
      studies: [
        'Volume@tv-basicstudies'
      ],
      overrides: theme === 'dark' ? {
        'paneProperties.background': '#1a1a1a',
        'paneProperties.vertGridProperties.color': '#2a2a2a',
        'paneProperties.horzGridProperties.color': '#2a2a2a',
        'symbolWatermarkProperties.transparency': 90,
        'scalesProperties.textColor': '#b0b0b0',
        'mainSeriesProperties.candleStyle.upColor': '#10b981',
        'mainSeriesProperties.candleStyle.downColor': '#ef4444',
        'mainSeriesProperties.candleStyle.borderUpColor': '#10b981',
        'mainSeriesProperties.candleStyle.borderDownColor': '#ef4444',
        'mainSeriesProperties.candleStyle.wickUpColor': '#10b981',
        'mainSeriesProperties.candleStyle.wickDownColor': '#ef4444'
      } : {},
      disabled_features: [
        'use_localstorage_for_settings',
        'volume_force_overlay',
        'header_symbol_search',
        'header_resolutions',
        'header_chart_type',
        'header_settings',
        'header_indicators',
        'header_compare',
        'header_undo_redo',
        'header_screenshot',
        'header_fullscreen_button'
      ],
      enabled_features: [
        'study_templates'
      ]
    });

    // Add script to container
    if (containerRef.current) {
      containerRef.current.appendChild(script);
      scriptRef.current = script;
    }

    // Cleanup function
    return () => {
      if (scriptRef.current && containerRef.current) {
        try {
          containerRef.current.removeChild(scriptRef.current);
        } catch (e) {
          // Script might have been removed already
        }
      }
    };
  }, [symbol, theme, interval, width, height, autosize]);

  return (
    <div className={`tradingview-widget-container ${className}`}>
      <div
        ref={containerRef}
        id={`tradingview-widget-${Math.random().toString(36).substr(2, 9)}`}
        className="tradingview-widget"
        style={{
          width: autosize ? '100%' : width,
          height: autosize ? '100%' : height,
          minHeight: '400px'
        }}
      />
      <div className="tradingview-widget-copyright">
        <a 
          href="https://www.tradingview.com/" 
          rel="noopener nofollow" 
          target="_blank"
          className="text-xs text-gray-500 hover:text-gray-400"
        >
          Track all markets on TradingView
        </a>
      </div>
    </div>
  );
});

TradingViewWidget.displayName = 'TradingViewWidget';

export default TradingViewWidget;
