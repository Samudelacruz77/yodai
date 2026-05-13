# YodAI

An AI chatbot that speaks like Yoda, running locally on NVIDIA Jetson AGX Orin with TensorRT-LLM.

## Architecture

```
Browser  <-->  Go HTTP Server (:8080)  <-->  TensorRT-LLM Server (:8000)
               (SSE streaming)                (Llama 3.1 8B INT4 AWQ)
               (Yoda system prompt)           (OpenAI-compatible API)
```

## Prerequisites

- NVIDIA Jetson AGX Orin with JetPack 6.1+ (L4T r36.3+)
- Docker with nvidia-container-runtime
- Docker Compose v2
- ~20 GB free disk space for model weights and engine
- Hugging Face account with Llama 3.1 access

## Quick Start

### 1. Build the TensorRT-LLM engine (one-time, ~20 min)

```bash
make engine HUGGINGFACE_TOKEN=hf_your_token_here
```

### 2. Build and start

```bash
make docker-build
make up
```

### 3. Open the chat

Navigate to `http://<jetson-ip>:8080` in your browser.

## Development

Run the Go server locally (requires inference backend running separately):

```bash
make dev
```

Build for the local architecture:

```bash
make build-local
```

Cross-compile for ARM64 (Jetson):

```bash
make build
```

Run tests:

```bash
make test
```

## Configuration

Configuration is loaded from (in order of precedence):
1. Environment variables (`YODAI_*` prefix)
2. YAML config file (`--config path/to/yodai.yaml`)
3. Defaults

| Variable | Default | Description |
|----------|---------|-------------|
| `YODAI_LISTEN_ADDR` | `:8080` | HTTP server listen address |
| `YODAI_INFERENCE_URL` | `http://localhost:8000` | TensorRT-LLM server URL |
| `YODAI_INFERENCE_MODEL` | `tensorrt_llm` | Model name for API requests |
| `YODAI_MAX_TOKENS` | `512` | Max tokens per response |
| `YODAI_TEMPERATURE` | `0.7` | Sampling temperature |
| `YODAI_TOP_P` | `0.9` | Nucleus sampling threshold |

## Project Structure

```
cmd/yodai/          Go entrypoint
internal/config/    Configuration loading
internal/chat/      Chat handlers, OpenAI client, Yoda prompt
internal/server/    HTTP router and middleware
web/                Embedded web UI (HTML/CSS/JS)
deploy/             Containerfiles, Compose, engine build scripts
configs/            Default configuration
```

## License

MIT
