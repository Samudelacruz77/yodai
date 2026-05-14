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
- Podman 4.1+ (or Docker 25+ with `CONTAINER_RUNTIME=docker`)
- NVIDIA Container Toolkit with CDI configured (`nvidia-ctk cdi generate --output=/etc/cdi/nvidia.yaml`)

## Quick Start

The inference image on quay.io includes a pre-compiled TensorRT-LLM engine — no on-device build needed.

### 1. Start the containers

```bash
make up
```

### 2. Open the chat

Navigate to `http://<jetson-ip>:8080` in your browser.

## Building the Engine (one-time, maintainer only)

If you need to rebuild the TensorRT-LLM engine (e.g. for a different model or Jetson variant):

### 1. Build the engine

```bash
make engine HUGGINGFACE_TOKEN=hf_your_token_here
```

### 2. Package and push the inference image

```bash
make inference-image
podman push quay.io/sdelacru/yodai-inference:latest
```

This exports the engine from the build volume, bakes it into the inference container image, and pushes to quay.io. Other devices just pull the image.

## Bootc Deployment

Deploy YodAI as part of an immutable RHEL bootc image for the Jetson.

### Build the bootc image

```bash
podman build -f deploy/bootc/Containerfile -t quay.io/sdelacru/yodai-bootc:latest .
podman push quay.io/sdelacru/yodai-bootc:latest
```

### Deploy to a Jetson

```bash
bootc switch quay.io/sdelacru/yodai-bootc:latest
systemctl reboot
```

The boot sequence:

1. **yodai-cdi-generate** — generates the NVIDIA CDI spec (skips if already done)
2. **yodai** — pulls container images and starts the compose stack

### Flightctl

A fleet spec is provided for managing Jetson devices with Flightctl:

```bash
flightctl apply -f deploy/flightctl/fleet.yaml
flightctl label device <device-name> fleet=yodai
```

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
cmd/yodai/              Go entrypoint
internal/config/        Configuration loading
internal/chat/          Chat handlers, OpenAI client, Yoda prompt
internal/server/        HTTP router and middleware
web/                    Embedded web UI (HTML/CSS/JS)
deploy/
  Containerfile.*       Application container images
  docker-compose.yml    Local development compose
  scripts/              Engine build and inference scripts
  bootc/                Bootc image (Containerfile, systemd services, compose)
  flightctl/            Flightctl fleet spec
configs/                Default configuration
```

## License

MIT
