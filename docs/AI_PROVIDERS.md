# AI Providers Integration Guide

This document provides comprehensive information about integrating and using different AI providers with the AI Agentic Crypto Browser.

## Supported Providers

The system supports the following AI providers:

1. **OpenAI** - Cloud-based GPT models
2. **Anthropic** - Cloud-based Claude models  
3. **Ollama** - Local LLM hosting platform
4. **LM Studio** - Local model management and inference

## Configuration

### Environment Variables

Set the following environment variables to configure AI providers:

```bash
# General AI Configuration
AI_MODEL_PROVIDER=ollama  # openai, anthropic, ollama, lmstudio

# OpenAI Configuration
OPENAI_API_KEY=your_openai_api_key
AI_MODEL_NAME=gpt-4-turbo-preview

# Anthropic Configuration  
ANTHROPIC_API_KEY=your_anthropic_api_key

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

### Configuration File

The AI configuration is also managed through `configs/ai.yaml`:

```yaml
providers:
  ollama:
    enabled: true
    base_url: "${OLLAMA_BASE_URL:-http://localhost:11434}"
    models:
      - name: "qwen3"
        temperature: 0.7
        top_p: 1.0
        top_k: 40
        num_ctx: 2048
      - name: "codellama"
        temperature: 0.3
        top_p: 0.9
        top_k: 40
        num_ctx: 4096
    timeout: "300s"
    max_retries: 3
    retry_delay: "2s"
    health_check_interval: "30s"

  lmstudio:
    enabled: true
    base_url: "${LMSTUDIO_BASE_URL:-http://localhost:1234/v1}"
    models:
      - name: "local-model"
        max_tokens: 4000
        temperature: 0.7
        top_p: 1.0
      - name: "llama-2-7b-chat"
        max_tokens: 4000
        temperature: 0.7
        top_p: 1.0
    timeout: "300s"
    max_retries: 3
    retry_delay: "2s"
    health_check_interval: "30s"
```

## Ollama Integration

### Prerequisites

1. Install Ollama from [https://ollama.ai](https://ollama.ai)
2. Pull desired models:
   ```bash
   ollama pull qwen3
   ollama pull codellama
   ollama pull mistral
   ```
3. Start Ollama service:
   ```bash
   ollama serve
   ```

### Supported Models

- **qwen3** - General purpose conversational model
- **qwen3:4b** - Larger variant with better performance
- **codellama** - Specialized for code generation and analysis
- **mistral** - Fast and efficient general purpose model
- **neural-chat** - Optimized for conversational AI

### Configuration Parameters

- **base_url**: Ollama server endpoint (default: http://localhost:11434)
- **model**: Default model to use
- **temperature**: Controls randomness (0.0-1.0)
- **top_p**: Nucleus sampling parameter
- **top_k**: Top-k sampling parameter
- **num_ctx**: Context window size
- **timeout**: Request timeout duration
- **max_retries**: Maximum retry attempts
- **retry_delay**: Delay between retries

### Usage Examples

```go
// Using Ollama provider
config := ai.OllamaConfig{
    BaseURL:     "http://localhost:11434",
    Model:       "qwen3",
    Temperature: 0.7,
    TopP:        1.0,
    TopK:        40,
    NumCtx:      2048,
    Timeout:     300 * time.Second,
}

provider := ai.NewOllamaProvider(config, logger)

// Generate response
messages := []ai.Message{
    {Role: "user", Content: "Explain quantum computing"},
}
response, err := provider.GenerateResponse(ctx, messages)
```

## LM Studio Integration

### Prerequisites

1. Download and install LM Studio from [https://lmstudio.ai](https://lmstudio.ai)
2. Load a model in LM Studio
3. Start the local server (usually on port 1234)

### Supported Models

LM Studio supports various model formats:
- **GGUF models** (recommended)
- **Llama 2** variants
- **Code Llama** models
- **Mistral** models
- **Custom fine-tuned models**

### Configuration Parameters

- **base_url**: LM Studio server endpoint (default: http://localhost:1234/v1)
- **model**: Model identifier in LM Studio
- **temperature**: Controls randomness (0.0-1.0)
- **max_tokens**: Maximum tokens in response
- **top_p**: Nucleus sampling parameter
- **timeout**: Request timeout duration
- **max_retries**: Maximum retry attempts
- **retry_delay**: Delay between retries

### Usage Examples

```go
// Using LM Studio provider
config := ai.LMStudioConfig{
    BaseURL:     "http://localhost:1234/v1",
    Model:       "llama-2-7b-chat",
    Temperature: 0.7,
    MaxTokens:   4000,
    TopP:        1.0,
    Timeout:     300 * time.Second,
}

provider := ai.NewLMStudioProvider(config, logger)

// Analyze content
analysis, err := provider.AnalyzeContent(ctx, content, "sentiment")
```

## Health Monitoring

### Health Check Endpoints

The system provides comprehensive health monitoring for all AI providers:

```bash
# Check overall AI health
GET /health/ai

# Check specific provider health
GET /health/ai/{provider}

# Trigger immediate health check
POST /health/ai/{provider}/check

# Get available models for provider
GET /health/ai/{provider}/models
```

### Health Status Response

```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "providers": {
    "ollama": {
      "provider": "ollama",
      "healthy": true,
      "last_checked": "2024-01-15T10:29:45Z",
      "error": "",
      "models": ["qwen3", "codellama", "mistral"],
      "response_time": "150ms"
    },
    "lmstudio": {
      "provider": "lmstudio", 
      "healthy": true,
      "last_checked": "2024-01-15T10:29:50Z",
      "error": "",
      "models": ["llama-2-7b-chat", "code-llama-7b-instruct"],
      "response_time": "200ms"
    }
  },
  "summary": {
    "total": 2,
    "healthy": 2,
    "unhealthy": 0
  }
}
```

### Monitoring Features

- **Automatic Health Checks**: Periodic health monitoring
- **Real-time Status**: Current provider availability
- **Model Discovery**: Automatic detection of available models
- **Response Time Tracking**: Performance monitoring
- **Error Reporting**: Detailed error information
- **Graceful Degradation**: Fallback to healthy providers

## API Usage

### Chat Completion

```bash
curl -X POST http://localhost:8082/ai/chat \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_jwt_token" \
  -d '{
    "message": "Help me analyze this website",
    "conversation_id": "optional-conversation-id"
  }'
```

### Content Analysis

```bash
curl -X POST http://localhost:8082/ai/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_jwt_token" \
  -d '{
    "type": "analyze_content",
    "input": {
      "url": "https://example.com",
      "analysis_type": "sentiment"
    }
  }'
```

## Troubleshooting

### Common Issues

1. **Connection Refused**
   - Ensure Ollama/LM Studio is running
   - Check port configuration
   - Verify firewall settings

2. **Model Not Found**
   - Pull/load the specified model
   - Check model name spelling
   - Verify model availability

3. **Timeout Errors**
   - Increase timeout configuration
   - Check system resources
   - Consider using smaller models

4. **Memory Issues**
   - Monitor system RAM usage
   - Use quantized models
   - Adjust context window size

### Debug Commands

```bash
# Check Ollama status
ollama list
ollama ps

# Test Ollama API
curl http://localhost:11434/api/tags

# Test LM Studio API  
curl http://localhost:1234/v1/models
```

### Logs

Monitor application logs for detailed error information:

```bash
# View AI service logs
docker logs ai-agentic-browser-ai-agent-1

# Check health monitoring
grep "health check" /var/log/ai-agent.log
```

## Performance Optimization

### Model Selection

- **Small models** (7B parameters): Faster inference, lower memory
- **Large models** (13B+ parameters): Better quality, higher resource usage
- **Specialized models**: Optimized for specific tasks

### Configuration Tuning

- **Context Window**: Balance between capability and performance
- **Temperature**: Lower for deterministic outputs, higher for creativity
- **Batch Size**: Optimize for throughput vs latency
- **Caching**: Enable model caching for faster startup

### Hardware Recommendations

- **CPU**: Multi-core processor (8+ cores recommended)
- **RAM**: 16GB+ for 7B models, 32GB+ for 13B models
- **GPU**: Optional but significantly improves performance
- **Storage**: SSD for faster model loading

## Security Considerations

### Local Deployment Benefits

- **Data Privacy**: All processing happens locally
- **No API Keys**: No cloud service dependencies
- **Network Isolation**: Can run completely offline
- **Compliance**: Easier regulatory compliance

### Best Practices

- **Access Control**: Secure API endpoints
- **Model Validation**: Verify model integrity
- **Resource Limits**: Prevent resource exhaustion
- **Monitoring**: Track usage and performance
- **Updates**: Keep software and models updated
