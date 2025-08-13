#!/bin/bash

# Firebase Setup Script for AI Agentic Crypto Browser
# This script initializes Firebase for the project

set -e

echo "🔥 Setting up Firebase for AI Agentic Crypto Browser..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if Firebase CLI is installed
if ! command -v firebase &> /dev/null; then
    echo -e "${RED}Firebase CLI is not installed. Installing...${NC}"
    npm install -g firebase-tools
fi

# Check if user is logged in to Firebase
if ! firebase projects:list &> /dev/null; then
    echo -e "${YELLOW}Please log in to Firebase...${NC}"
    firebase login
fi

# Project configuration
PROJECT_ID="ai-agentic-crypto-browser"
PROJECT_NAME="AI Agentic Crypto Browser"
REGION="us-central1"

echo -e "${BLUE}Setting up Firebase project: ${PROJECT_ID}${NC}"

# Initialize Firebase in the current directory
if [ ! -f "firebase.json" ]; then
    echo -e "${YELLOW}Initializing Firebase project...${NC}"
    
    # Create firebase.json configuration
    cat > firebase.json << EOF
{
  "firestore": {
    "rules": "firestore.rules",
    "indexes": "firestore.indexes.json"
  },
  "database": {
    "rules": "database.rules.json"
  },
  "storage": {
    "rules": "storage.rules"
  },
  "functions": {
    "source": "functions",
    "runtime": "nodejs18"
  },
  "hosting": {
    "public": "web/dist",
    "ignore": [
      "firebase.json",
      "**/.*",
      "**/node_modules/**"
    ],
    "rewrites": [
      {
        "source": "/api/**",
        "function": "api"
      },
      {
        "source": "**",
        "destination": "/index.html"
      }
    ]
  },
  "emulators": {
    "auth": {
      "port": 9099
    },
    "firestore": {
      "port": 8080
    },
    "database": {
      "port": 9000
    },
    "storage": {
      "port": 9199
    },
    "functions": {
      "port": 5001
    },
    "hosting": {
      "port": 5000
    },
    "pubsub": {
      "port": 8085
    },
    "ui": {
      "enabled": true,
      "port": 4000
    }
  }
}
EOF

    echo -e "${GREEN}Created firebase.json${NC}"
fi

# Create Firestore security rules
if [ ! -f "firestore.rules" ]; then
    cat > firestore.rules << 'EOF'
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    // Users can read and write their own data
    match /user_portfolios/{userId} {
      allow read, write: if request.auth != null && request.auth.uid == userId;
    }
    
    match /user_preferences/{userId} {
      allow read, write: if request.auth != null && request.auth.uid == userId;
    }
    
    match /user_notifications/{userId} {
      allow read, write: if request.auth != null && request.auth.uid == userId;
    }
    
    // Trading signals - read for authenticated users, write for admins
    match /trading_signals/{signalId} {
      allow read: if request.auth != null;
      allow write: if request.auth != null && 
        ('admin' in request.auth.token.roles || 'trader' in request.auth.token.roles);
    }
    
    // Market data - read for authenticated users
    match /market_data/{dataId} {
      allow read: if request.auth != null;
      allow write: if request.auth != null && 'admin' in request.auth.token.roles;
    }
    
    // Trading strategies - users can manage their own
    match /trading_strategies/{strategyId} {
      allow read, write: if request.auth != null && 
        (resource.data.user_id == request.auth.uid || 'admin' in request.auth.token.roles);
    }
    
    // Risk metrics - read for authenticated users
    match /risk_metrics/{metricId} {
      allow read: if request.auth != null;
      allow write: if request.auth != null && 'admin' in request.auth.token.roles;
    }
    
    // Performance analytics - read for authenticated users
    match /performance_analytics/{analyticsId} {
      allow read: if request.auth != null;
      allow write: if request.auth != null && 'admin' in request.auth.token.roles;
    }
    
    // System logs - admin only
    match /system_logs/{logId} {
      allow read, write: if request.auth != null && 'admin' in request.auth.token.roles;
    }
    
    // Audit trail - admin only
    match /audit_trail/{auditId} {
      allow read, write: if request.auth != null && 'admin' in request.auth.token.roles;
    }
    
    // Default deny
    match /{document=**} {
      allow read, write: if false;
    }
  }
}
EOF
    echo -e "${GREEN}Created firestore.rules${NC}"
fi

# Create Firestore indexes
if [ ! -f "firestore.indexes.json" ]; then
    cat > firestore.indexes.json << 'EOF'
{
  "indexes": [
    {
      "collectionGroup": "trading_signals",
      "queryScope": "COLLECTION",
      "fields": [
        {
          "fieldPath": "symbol",
          "order": "ASCENDING"
        },
        {
          "fieldPath": "timestamp",
          "order": "DESCENDING"
        }
      ]
    },
    {
      "collectionGroup": "market_data",
      "queryScope": "COLLECTION",
      "fields": [
        {
          "fieldPath": "symbol",
          "order": "ASCENDING"
        },
        {
          "fieldPath": "timestamp",
          "order": "DESCENDING"
        },
        {
          "fieldPath": "exchange",
          "order": "ASCENDING"
        }
      ]
    },
    {
      "collectionGroup": "user_portfolios",
      "queryScope": "COLLECTION",
      "fields": [
        {
          "fieldPath": "user_id",
          "order": "ASCENDING"
        },
        {
          "fieldPath": "created_at",
          "order": "DESCENDING"
        }
      ]
    },
    {
      "collectionGroup": "trading_strategies",
      "queryScope": "COLLECTION",
      "fields": [
        {
          "fieldPath": "user_id",
          "order": "ASCENDING"
        },
        {
          "fieldPath": "status",
          "order": "ASCENDING"
        },
        {
          "fieldPath": "created_at",
          "order": "DESCENDING"
        }
      ]
    }
  ],
  "fieldOverrides": []
}
EOF
    echo -e "${GREEN}Created firestore.indexes.json${NC}"
fi

# Create Realtime Database rules
if [ ! -f "database.rules.json" ]; then
    cat > database.rules.json << 'EOF'
{
  "rules": {
    "live_prices": {
      ".read": "auth != null",
      ".write": "auth != null && auth.token.roles.includes('admin')"
    },
    "active_orders": {
      "$uid": {
        ".read": "auth != null && auth.uid == $uid",
        ".write": "auth != null && auth.uid == $uid"
      }
    },
    "user_sessions": {
      "$uid": {
        ".read": "auth != null && auth.uid == $uid",
        ".write": "auth != null && auth.uid == $uid"
      }
    },
    "system_status": {
      ".read": "auth != null",
      ".write": "auth != null && auth.token.roles.includes('admin')"
    },
    "alerts": {
      "$uid": {
        ".read": "auth != null && auth.uid == $uid",
        ".write": "auth != null && auth.uid == $uid"
      }
    },
    "chat_rooms": {
      ".read": "auth != null",
      ".write": "auth != null"
    }
  }
}
EOF
    echo -e "${GREEN}Created database.rules.json${NC}"
fi

# Create Storage rules
if [ ! -f "storage.rules" ]; then
    cat > storage.rules << 'EOF'
rules_version = '2';
service firebase.storage {
  match /b/{bucket}/o {
    // Users can upload to their own folder
    match /users/{userId}/{allPaths=**} {
      allow read, write: if request.auth != null && request.auth.uid == userId;
    }
    
    // Public read access for certain files
    match /public/{allPaths=**} {
      allow read: if true;
      allow write: if request.auth != null && 'admin' in request.auth.token.roles;
    }
    
    // Trading reports - authenticated users can read, admins can write
    match /reports/{allPaths=**} {
      allow read: if request.auth != null;
      allow write: if request.auth != null && 'admin' in request.auth.token.roles;
    }
    
    // System backups - admin only
    match /backups/{allPaths=**} {
      allow read, write: if request.auth != null && 'admin' in request.auth.token.roles;
    }
  }
}
EOF
    echo -e "${GREEN}Created storage.rules${NC}"
fi

# Set up environment variables
echo -e "${YELLOW}Setting up environment variables...${NC}"

# Add Firebase configuration to .env file
if [ ! -f ".env" ]; then
    touch .env
fi

# Check if Firebase config already exists in .env
if ! grep -q "FIREBASE_PROJECT_ID" .env; then
    cat >> .env << EOF

# Firebase Configuration
FIREBASE_PROJECT_ID=${PROJECT_ID}
FIREBASE_DATABASE_URL=https://${PROJECT_ID}-default-rtdb.firebaseio.com
FIREBASE_STORAGE_BUCKET=${PROJECT_ID}.appspot.com
FIREBASE_CREDENTIALS_PATH=./configs/firebase-service-account.json

# Firebase Emulator Configuration (for development)
FIREBASE_AUTH_EMULATOR_HOST=localhost:9099
FIRESTORE_EMULATOR_HOST=localhost:8080
FIREBASE_DATABASE_EMULATOR_HOST=localhost:9000
FIREBASE_STORAGE_EMULATOR_HOST=localhost:9199
FUNCTIONS_EMULATOR_HOST=localhost:5001
EOF
    echo -e "${GREEN}Added Firebase configuration to .env${NC}"
fi

# Create a sample service account key file template
if [ ! -f "configs/firebase-service-account.json.template" ]; then
    mkdir -p configs
    cat > configs/firebase-service-account.json.template << 'EOF'
{
  "type": "service_account",
  "project_id": "your-project-id",
  "private_key_id": "your-private-key-id",
  "private_key": "-----BEGIN PRIVATE KEY-----\nYOUR_PRIVATE_KEY\n-----END PRIVATE KEY-----\n",
  "client_email": "your-service-account@your-project-id.iam.gserviceaccount.com",
  "client_id": "your-client-id",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/your-service-account%40your-project-id.iam.gserviceaccount.com"
}
EOF
    echo -e "${GREEN}Created service account template at configs/firebase-service-account.json.template${NC}"
fi

echo -e "${GREEN}✅ Firebase setup completed!${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Create a Firebase project at https://console.firebase.google.com/"
echo "2. Download the service account key and save it as configs/firebase-service-account.json"
echo "3. Update the project ID in your configuration files"
echo "4. Run 'firebase deploy' to deploy your rules and indexes"
echo "5. For local development, run 'firebase emulators:start'"
echo ""
echo -e "${BLUE}Useful commands:${NC}"
echo "• firebase emulators:start - Start Firebase emulators for local development"
echo "• firebase deploy - Deploy to production"
echo "• firebase deploy --only firestore:rules - Deploy only Firestore rules"
echo "• firebase deploy --only firestore:indexes - Deploy only Firestore indexes"
echo "• firebase projects:list - List your Firebase projects"
echo ""
echo -e "${GREEN}Happy coding! 🚀${NC}"
