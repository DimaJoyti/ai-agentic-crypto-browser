# Ollama and LM Studio Setup Guide

This guide walks you through setting up Ollama and LM Studio integrations with the AI Agentic Crypto Browser.

## Quick Start

### Option 1: Ollama Setup (Recommended for beginners)

1. **Install Ollama**
   ```bash
   # macOS
   brew install ollama
   
   # Linux
   curl -fsSL https://ollama.ai/install.sh | sh
   
   # Windows
   # Download from https://ollama.ai/download
   ```

2. **Start Ollama Service**
   ```bash
   ollama serve
   ```

3. **Pull a Model**
   ```bash
   # Start with a smaller model
   ollama pull qwen3:4b
   
   # Or for better performance (requires more RAM)
   ollama pull qwen3:4b
   ```

4. **Configure Environment**
   ```bash
   export AI_MODEL_PROVIDER=ollama
   export OLLAMA_MODEL=qwen3:4b
   ```

5. **Start the Application**
   ```bash
   docker-compose up -d
   ```

### Option 2: LM Studio Setup

1. **Download LM Studio**
   - Visit [https://lmstudio.ai](https://lmstudio.ai)
   - Download for your platform (Windows, macOS, Linux)

2. **Install and Launch LM Studio**
   - Follow the installation wizard
   - Launch the application

3. **Download a Model**
   - Browse the model library
   - Download a model (e.g., "Llama 2 7B Chat")
   - Wait for download to complete

4. **Load the Model**
   - Go to the "Chat" tab
   - Select your downloaded model
   - Click "Load Model"

5. **Start Local Server**
   - Go to "Local Server" tab
   - Click "Start Server"
   - Note the server URL (usually http://localhost:1234)

6. **Configure Environment**
   ```bash
   export AI_MODEL_PROVIDER=lmstudio
   export LMSTUDIO_MODEL=your-model-name
   ```

7. **Start the Application**
   ```bash
   docker-compose up -d
   ```

## Detailed Setup Instructions

### Ollama Detailed Setup

#### 1. Installation

**macOS:**
```bash
# Using Homebrew
brew install ollama

# Or download from website
curl -L https://ollama.ai/download/ollama-darwin.zip -o ollama.zip
unzip ollama.zip
sudo mv ollama /usr/local/bin/
```

**Linux:**
```bash
# Automatic installation
curl -fsSL https://ollama.ai/install.sh | sh

# Manual installation
wget https://ollama.ai/download/ollama-linux-amd64.tgz
tar -xzf ollama-linux-amd64.tgz
sudo mv ollama /usr/local/bin/
```

**Windows:**
- Download the installer from [https://ollama.ai/download](https://ollama.ai/download)
- Run the installer as administrator
- Add Ollama to your PATH

#### 2. Model Management

**Available Models:**
```bash
# List available models online
ollama list

# Pull specific models
ollama pull qwen3          # ~3.8GB
ollama pull qwen3:4b      # ~7.3GB  
ollama pull codellama       # ~3.8GB
ollama pull mistral         # ~4.1GB
ollama pull neural-chat     # ~4.1GB

# Check downloaded models
ollama list
```

**Model Recommendations:**
- **Beginners**: `qwen3:4b` (3.8GB, good balance)
- **Code Tasks**: `codellama` (3.8GB, specialized for coding)
- **Performance**: `qwen3:4b` (7.3GB, better quality)
- **Speed**: `mistral` (4.1GB, faster inference)

#### 3. Service Configuration

**Start as Service (Linux/macOS):**
```bash
# Create systemd service
sudo tee /etc/systemd/system/ollama.service > /dev/null <<EOF
[Unit]
Description=Ollama Service
After=network.target

[Service]
Type=simple
User=ollama
Group=ollama
ExecStart=/usr/local/bin/ollama serve
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
EOF

# Enable and start service
sudo systemctl enable ollama
sudo systemctl start ollama
```

**Manual Start:**
```bash
# Start in foreground
ollama serve

# Start in background
nohup ollama serve > ollama.log 2>&1 &
```

#### 4. Testing Ollama

```bash
# Test API directly
curl http://localhost:11434/api/tags

# Test chat completion
curl http://localhost:11434/api/chat -d '{
  "model": "qwen3",
  "messages": [
    {"role": "user", "content": "Hello!"}
  ]
}'
```

### LM Studio Detailed Setup

#### 1. Installation

**Download and Install:**
1. Visit [https://lmstudio.ai](https://lmstudio.ai)
2. Click "Download" for your platform
3. Run the installer
4. Launch LM Studio

#### 2. Model Management

**Downloading Models:**
1. Open LM Studio
2. Go to "Discover" tab
3. Browse available models
4. Click "Download" on desired model
5. Wait for download (models are 3-20GB)

**Recommended Models:**
- **Llama 2 7B Chat** - General purpose, good balance
- **Code Llama 7B Instruct** - Code generation and analysis
- **Mistral 7B Instruct** - Fast and efficient
- **Neural Chat 7B** - Conversational AI optimized

#### 3. Loading Models

**Load a Model:**
1. Go to "Chat" tab
2. Click model selector dropdown
3. Choose downloaded model
4. Click "Load Model"
5. Wait for model to load into memory

**Model Settings:**
- **Context Length**: Adjust based on use case
- **GPU Acceleration**: Enable if available
- **Temperature**: Control randomness (0.1-1.0)
- **Max Tokens**: Limit response length

#### 4. Starting Local Server

**Enable API Server:**
1. Go to "Local Server" tab
2. Select loaded model
3. Configure server settings:
   - Port: 1234 (default)
   - CORS: Enable for web access
   - API Key: Optional
4. Click "Start Server"

**Server Configuration:**
```json
{
  "host": "localhost",
  "port": 1234,
  "cors": true,
  "api_key": "",
  "max_requests_per_minute": 60
}
```

#### 5. Testing LM Studio

```bash
# Test server status
curl http://localhost:1234/v1/models

# Test chat completion
curl http://localhost:1234/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "local-model",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ]
  }'
```

## Application Configuration

### Environment Variables

Create a `.env` file in your project root:

```bash
# Choose your provider
AI_MODEL_PROVIDER=ollama  # or lmstudio

# Ollama Configuration
OLLAMA_BASE_URL=http://localhost:11434
OLLAMA_MODEL=qwen3
OLLAMA_TEMPERATURE=0.7
OLLAMA_TOP_P=1.0
OLLAMA_TOP_K=40
OLLAMA_NUM_CTX=2048
OLLAMA_TIMEOUT=300s
OLLAMA_MAX_RETRIES=3
OLLAMA_RETRY_DELAY=2s
OLLAMA_HEALTH_CHECK_INTERVAL=30s

# LM Studio Configuration  
LMSTUDIO_BASE_URL=http://localhost:1234/v1
LMSTUDIO_MODEL=local-model
LMSTUDIO_TEMPERATURE=0.7
LMSTUDIO_MAX_TOKENS=4000
LMSTUDIO_TOP_P=1.0
LMSTUDIO_TIMEOUT=300s
LMSTUDIO_MAX_RETRIES=3
LMSTUDIO_RETRY_DELAY=2s
LMSTUDIO_HEALTH_CHECK_INTERVAL=30s
```

### Docker Compose

Update your `docker-compose.yml`:

```yaml
version: '3.8'
services:
  ai-agent:
    build: ./cmd/ai-agent
    environment:
      - AI_MODEL_PROVIDER=${AI_MODEL_PROVIDER}
      - OLLAMA_BASE_URL=${OLLAMA_BASE_URL}
      - OLLAMA_MODEL=${OLLAMA_MODEL}
      - LMSTUDIO_BASE_URL=${LMSTUDIO_BASE_URL}
      - LMSTUDIO_MODEL=${LMSTUDIO_MODEL}
    ports:
      - "8082:8082"
    depends_on:
      - postgres
      - redis
    networks:
      - app-network
    # For Ollama/LM Studio access
    extra_hosts:
      - "host.docker.internal:host-gateway"
```

## Verification and Testing

### Health Checks

```bash
# Check overall AI health
curl http://localhost:8082/health/ai

# Check Ollama specifically
curl http://localhost:8082/health/ai/ollama

# Check LM Studio specifically  
curl http://localhost:8082/health/ai/lmstudio

# Get available models
curl http://localhost:8082/health/ai/ollama/models
```

### Test Chat

```bash
# Get JWT token first (see authentication docs)
export JWT_TOKEN="your_jwt_token"

# Test chat with Ollama/LM Studio
curl -X POST http://localhost:8082/ai/chat \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "message": "Hello! Can you help me analyze a website?"
  }'
```

### Performance Testing

```bash
# Test response time
time curl -X POST http://localhost:8082/ai/chat \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{"message": "What is quantum computing?"}'

# Monitor resource usage
htop  # or Activity Monitor on macOS
```

## Troubleshooting

### Common Issues

**1. Connection Refused**
```bash
# Check if Ollama is running
ps aux | grep ollama
curl http://localhost:11434/api/tags

# Check if LM Studio server is running
curl http://localhost:1234/v1/models
```

**2. Model Not Found**
```bash
# List Ollama models
ollama list

# Pull missing model
ollama pull qwen3
```

**3. Out of Memory**
```bash
# Check memory usage
free -h  # Linux
vm_stat  # macOS

# Use smaller model
ollama pull qwen3:4b  # instead of 13b
```

**4. Slow Performance**
- Close other applications
- Use GPU acceleration if available
- Reduce context window size
- Use quantized models

### Debug Commands

```bash
# Ollama logs
ollama logs

# Check Ollama processes
ollama ps

# Test model directly
ollama run qwen3 "Hello, world!"

# LM Studio logs
# Check LM Studio console/logs tab
```

### Getting Help

- **Ollama**: [https://github.com/ollama/ollama](https://github.com/ollama/ollama)
- **LM Studio**: [https://lmstudio.ai/docs](https://lmstudio.ai/docs)
- **Project Issues**: [GitHub Issues](https://github.com/your-repo/issues)

## Next Steps

1. **Explore Models**: Try different models for various tasks
2. **Optimize Performance**: Tune parameters for your use case
3. **Monitor Usage**: Set up monitoring and alerting
4. **Scale Up**: Consider multiple instances or GPU acceleration
5. **Custom Models**: Train or fine-tune models for specific needs

## Advanced Configuration

### Multiple Providers

You can configure multiple providers and switch between them:

```bash
# Set primary provider
export AI_MODEL_PROVIDER=ollama

# Configure fallback in application code
# The system will automatically fall back to healthy providers
```

### Custom Models

**Ollama Custom Models:**
```bash
# Create Modelfile
echo 'FROM qwen3
PARAMETER temperature 0.8
SYSTEM "You are a helpful assistant."' > Modelfile

# Build custom model
ollama create my-custom-model -f Modelfile
```

**LM Studio Custom Models:**
1. Download GGUF model files
2. Place in LM Studio models directory
3. Refresh model list in LM Studio
4. Load and use custom model

This completes the setup guide for Ollama and LM Studio integrations!
