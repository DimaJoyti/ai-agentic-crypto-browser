/**
 * Chunk Loading Error Handler
 * Handles ChunkLoadError by automatically retrying or reloading the page
 */

interface ChunkLoadError extends Error {
  name: 'ChunkLoadError';
  message: string;
}

class ChunkErrorHandler {
  private retryCount = 0;
  private maxRetries = 3;
  private retryDelay = 1000; // 1 second

  constructor() {
    this.setupErrorHandlers();
  }

  private setupErrorHandlers() {
    // Handle unhandled promise rejections (chunk loading failures)
    if (typeof window !== 'undefined') {
      window.addEventListener('unhandledrejection', this.handleUnhandledRejection.bind(this));
      
      // Handle general errors
      window.addEventListener('error', this.handleError.bind(this));
      
      // Override webpack's chunk loading error handler
      this.setupWebpackErrorHandler();
    }
  }

  private handleUnhandledRejection(event: PromiseRejectionEvent) {
    const error = event.reason;
    
    if (this.isChunkLoadError(error)) {
      console.warn('Chunk loading error detected:', error);
      event.preventDefault(); // Prevent the error from being logged to console
      this.handleChunkError(error);
    }
  }

  private handleError(event: ErrorEvent) {
    if (this.isChunkLoadError(event.error)) {
      console.warn('Chunk loading error detected via error event:', event.error);
      event.preventDefault();
      this.handleChunkError(event.error);
    }
  }

  private isChunkLoadError(error: any): error is ChunkLoadError {
    return (
      error &&
      (error.name === 'ChunkLoadError' ||
        (error.message && error.message.includes('Loading chunk')) ||
        (error.message && error.message.includes('timeout')))
    );
  }

  private async handleChunkError(error: ChunkLoadError) {
    console.log(`Handling chunk error (attempt ${this.retryCount + 1}/${this.maxRetries}):`, error.message);

    if (this.retryCount < this.maxRetries) {
      this.retryCount++;
      
      try {
        // Wait before retrying
        await this.delay(this.retryDelay * this.retryCount);
        
        // Try to reload the failed chunk by refreshing the page
        console.log('Attempting to recover from chunk error...');
        
        // Clear any cached chunks
        this.clearChunkCache();
        
        // Reload the page to get fresh chunks
        window.location.reload();
      } catch (retryError) {
        console.error('Failed to recover from chunk error:', retryError);
        
        if (this.retryCount >= this.maxRetries) {
          this.showUserError();
        }
      }
    } else {
      console.error('Max retries exceeded for chunk loading');
      this.showUserError();
    }
  }

  private clearChunkCache() {
    try {
      // Clear any webpack chunk cache
      if ('caches' in window) {
        caches.keys().then(names => {
          names.forEach(name => {
            if (name.includes('webpack') || name.includes('chunk')) {
              caches.delete(name);
            }
          });
        });
      }

      // Clear localStorage entries related to chunks
      Object.keys(localStorage).forEach(key => {
        if (key.includes('webpack') || key.includes('chunk')) {
          localStorage.removeItem(key);
        }
      });
    } catch (error) {
      console.warn('Failed to clear chunk cache:', error);
    }
  }

  private setupWebpackErrorHandler() {
    // Override webpack's default chunk loading error handler
    if (typeof window !== 'undefined' && (window as any).__webpack_require__) {
      const originalRequire = (window as any).__webpack_require__;
      
      (window as any).__webpack_require__ = function(moduleId: string) {
        try {
          return originalRequire(moduleId);
        } catch (error: any) {
          if (error && error.message && error.message.includes('Loading chunk')) {
            console.warn('Webpack chunk loading intercepted:', error);
            // Let our handler deal with it
            throw error;
          }
          throw error;
        }
      };
    }
  }

  private delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  private showUserError() {
    // Show a user-friendly error message
    const errorDiv = document.createElement('div');
    errorDiv.innerHTML = `
      <div style="
        position: fixed;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        background: #f8f9fa;
        border: 1px solid #dee2e6;
        border-radius: 8px;
        padding: 20px;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        z-index: 10000;
        max-width: 400px;
        text-align: center;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      ">
        <h3 style="margin: 0 0 10px 0; color: #495057;">Loading Error</h3>
        <p style="margin: 0 0 15px 0; color: #6c757d;">
          There was a problem loading the application. Please refresh the page to continue.
        </p>
        <button onclick="window.location.reload()" style="
          background: #007bff;
          color: white;
          border: none;
          padding: 8px 16px;
          border-radius: 4px;
          cursor: pointer;
          font-size: 14px;
        ">
          Refresh Page
        </button>
      </div>
    `;
    
    document.body.appendChild(errorDiv);
  }

  // Public method to reset retry count
  public reset() {
    this.retryCount = 0;
  }
}

// Create and export a singleton instance
export const chunkErrorHandler = new ChunkErrorHandler();

// Auto-initialize in browser environment
if (typeof window !== 'undefined') {
  // Reset retry count on successful navigation
  window.addEventListener('beforeunload', () => {
    chunkErrorHandler.reset();
  });
}
