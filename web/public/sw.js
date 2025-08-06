const CACHE_NAME = 'ai-browser-v1.0.0'
const STATIC_CACHE = 'ai-browser-static-v1.0.0'
const DYNAMIC_CACHE = 'ai-browser-dynamic-v1.0.0'

// Assets to cache immediately
const STATIC_ASSETS = [
  '/',
  '/trading',
  '/dashboard',
  '/web3',
  '/analytics',
  '/manifest.json',
  '/icons/icon-192x192.png',
  '/icons/icon-512x512.png'
]

// API endpoints to cache
const API_CACHE_PATTERNS = [
  /\/api\/market-data/,
  /\/api\/portfolio/,
  /\/api\/user-preferences/
]

// Install event - cache static assets
self.addEventListener('install', (event) => {
  console.log('Service Worker: Installing...')
  
  event.waitUntil(
    caches.open(STATIC_CACHE)
      .then((cache) => {
        console.log('Service Worker: Caching static assets')
        return cache.addAll(STATIC_ASSETS)
      })
      .then(() => {
        console.log('Service Worker: Static assets cached')
        return self.skipWaiting()
      })
      .catch((error) => {
        console.error('Service Worker: Failed to cache static assets', error)
      })
  )
})

// Activate event - clean up old caches
self.addEventListener('activate', (event) => {
  console.log('Service Worker: Activating...')
  
  event.waitUntil(
    caches.keys()
      .then((cacheNames) => {
        return Promise.all(
          cacheNames.map((cacheName) => {
            if (cacheName !== STATIC_CACHE && cacheName !== DYNAMIC_CACHE) {
              console.log('Service Worker: Deleting old cache', cacheName)
              return caches.delete(cacheName)
            }
          })
        )
      })
      .then(() => {
        console.log('Service Worker: Activated')
        return self.clients.claim()
      })
  )
})

// Fetch event - serve from cache or network
self.addEventListener('fetch', (event) => {
  const { request } = event
  const url = new URL(request.url)

  // Skip non-GET requests
  if (request.method !== 'GET') {
    return
  }

  // Skip chrome-extension requests
  if (url.protocol === 'chrome-extension:') {
    return
  }

  // Handle different types of requests
  if (request.destination === 'document') {
    // HTML pages - Network first, fallback to cache
    event.respondWith(handlePageRequest(request))
  } else if (isAPIRequest(request)) {
    // API requests - Cache first for specific endpoints
    event.respondWith(handleAPIRequest(request))
  } else if (isStaticAsset(request)) {
    // Static assets - Cache first
    event.respondWith(handleStaticRequest(request))
  } else {
    // Other requests - Network first
    event.respondWith(handleNetworkFirst(request))
  }
})

// Handle page requests (HTML)
async function handlePageRequest(request) {
  try {
    // Try network first
    const networkResponse = await fetch(request)
    
    // Cache successful responses
    if (networkResponse.ok) {
      const cache = await caches.open(DYNAMIC_CACHE)
      cache.put(request, networkResponse.clone())
    }
    
    return networkResponse
  } catch (error) {
    console.log('Service Worker: Network failed, trying cache for page')
    
    // Fallback to cache
    const cachedResponse = await caches.match(request)
    if (cachedResponse) {
      return cachedResponse
    }
    
    // Fallback to offline page
    return caches.match('/')
  }
}

// Handle API requests
async function handleAPIRequest(request) {
  const url = new URL(request.url)
  
  // Check if this API should be cached
  const shouldCache = API_CACHE_PATTERNS.some(pattern => pattern.test(url.pathname))
  
  if (shouldCache) {
    try {
      // Try cache first for cacheable APIs
      const cachedResponse = await caches.match(request)
      if (cachedResponse) {
        // Return cached response and update in background
        updateCacheInBackground(request)
        return cachedResponse
      }
      
      // Fetch from network and cache
      const networkResponse = await fetch(request)
      if (networkResponse.ok) {
        const cache = await caches.open(DYNAMIC_CACHE)
        cache.put(request, networkResponse.clone())
      }
      
      return networkResponse
    } catch (error) {
      console.log('Service Worker: API request failed', error)
      
      // Return cached version if available
      const cachedResponse = await caches.match(request)
      if (cachedResponse) {
        return cachedResponse
      }
      
      // Return offline response
      return new Response(
        JSON.stringify({ 
          error: 'Offline', 
          message: 'This feature is not available offline' 
        }),
        {
          status: 503,
          headers: { 'Content-Type': 'application/json' }
        }
      )
    }
  } else {
    // Non-cacheable APIs - network only
    return fetch(request)
  }
}

// Handle static assets
async function handleStaticRequest(request) {
  try {
    // Try cache first
    const cachedResponse = await caches.match(request)
    if (cachedResponse) {
      return cachedResponse
    }
    
    // Fetch from network and cache
    const networkResponse = await fetch(request)
    if (networkResponse.ok) {
      const cache = await caches.open(STATIC_CACHE)
      cache.put(request, networkResponse.clone())
    }
    
    return networkResponse
  } catch (error) {
    console.log('Service Worker: Static asset failed to load', error)
    
    // Return cached version if available
    const cachedResponse = await caches.match(request)
    return cachedResponse || new Response('Asset not available offline', { status: 404 })
  }
}

// Handle other requests with network first strategy
async function handleNetworkFirst(request) {
  try {
    const networkResponse = await fetch(request)
    
    // Cache successful responses
    if (networkResponse.ok) {
      const cache = await caches.open(DYNAMIC_CACHE)
      cache.put(request, networkResponse.clone())
    }
    
    return networkResponse
  } catch (error) {
    // Fallback to cache
    const cachedResponse = await caches.match(request)
    return cachedResponse || new Response('Resource not available offline', { status: 404 })
  }
}

// Update cache in background
async function updateCacheInBackground(request) {
  try {
    const networkResponse = await fetch(request)
    if (networkResponse.ok) {
      const cache = await caches.open(DYNAMIC_CACHE)
      cache.put(request, networkResponse.clone())
    }
  } catch (error) {
    console.log('Service Worker: Background cache update failed', error)
  }
}

// Helper functions
function isAPIRequest(request) {
  const url = new URL(request.url)
  return url.pathname.startsWith('/api/')
}

function isStaticAsset(request) {
  const url = new URL(request.url)
  return /\.(js|css|png|jpg|jpeg|gif|svg|woff|woff2|ttf|eot|ico)$/.test(url.pathname)
}

// Background sync for offline actions
self.addEventListener('sync', (event) => {
  console.log('Service Worker: Background sync triggered', event.tag)
  
  if (event.tag === 'trading-orders') {
    event.waitUntil(syncTradingOrders())
  } else if (event.tag === 'user-preferences') {
    event.waitUntil(syncUserPreferences())
  }
})

// Sync trading orders when back online
async function syncTradingOrders() {
  try {
    // Get pending orders from IndexedDB
    const pendingOrders = await getPendingOrders()
    
    for (const order of pendingOrders) {
      try {
        const response = await fetch('/api/trading/orders', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(order)
        })
        
        if (response.ok) {
          await removePendingOrder(order.id)
          console.log('Service Worker: Order synced successfully', order.id)
        }
      } catch (error) {
        console.error('Service Worker: Failed to sync order', order.id, error)
      }
    }
  } catch (error) {
    console.error('Service Worker: Background sync failed', error)
  }
}

// Sync user preferences when back online
async function syncUserPreferences() {
  try {
    const pendingPreferences = await getPendingPreferences()
    
    if (pendingPreferences) {
      const response = await fetch('/api/user/preferences', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(pendingPreferences)
      })
      
      if (response.ok) {
        await clearPendingPreferences()
        console.log('Service Worker: Preferences synced successfully')
      }
    }
  } catch (error) {
    console.error('Service Worker: Failed to sync preferences', error)
  }
}

// Push notification handler
self.addEventListener('push', (event) => {
  console.log('Service Worker: Push notification received')
  
  const options = {
    body: 'You have new trading alerts',
    icon: '/icons/icon-192x192.png',
    badge: '/icons/badge-72x72.png',
    vibrate: [200, 100, 200],
    data: {
      url: '/trading'
    },
    actions: [
      {
        action: 'view',
        title: 'View Details',
        icon: '/icons/view-action.png'
      },
      {
        action: 'dismiss',
        title: 'Dismiss',
        icon: '/icons/dismiss-action.png'
      }
    ]
  }
  
  if (event.data) {
    const data = event.data.json()
    options.body = data.message || options.body
    options.data = { ...options.data, ...data }
  }
  
  event.waitUntil(
    self.registration.showNotification('AI Browser', options)
  )
})

// Notification click handler
self.addEventListener('notificationclick', (event) => {
  console.log('Service Worker: Notification clicked', event.action)
  
  event.notification.close()
  
  if (event.action === 'view') {
    const url = event.notification.data?.url || '/'
    event.waitUntil(
      clients.openWindow(url)
    )
  }
})

// Placeholder functions for IndexedDB operations
async function getPendingOrders() {
  // Implementation would use IndexedDB to get pending orders
  return []
}

async function removePendingOrder(orderId) {
  // Implementation would remove order from IndexedDB
  console.log('Removing pending order:', orderId)
}

async function getPendingPreferences() {
  // Implementation would get pending preferences from IndexedDB
  return null
}

async function clearPendingPreferences() {
  // Implementation would clear pending preferences from IndexedDB
  console.log('Clearing pending preferences')
}
