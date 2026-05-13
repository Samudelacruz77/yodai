.PHONY: build run test clean engine container-build up down dev tidy

BINARY := bin/yodai
GOARCH ?= arm64
CONTAINER_RUNTIME ?= podman

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build -o $(BINARY) ./cmd/yodai/

build-local:
	go build -o $(BINARY) ./cmd/yodai/

run: build-local
	./$(BINARY)

dev:
	YODAI_INFERENCE_URL=http://localhost:8000 go run ./cmd/yodai/

test:
	go test ./...

tidy:
	go mod tidy

container-build:
	$(CONTAINER_RUNTIME) compose -f deploy/docker-compose.yml build

engine:
	@if [ -z "$(HUGGINGFACE_TOKEN)" ]; then \
		echo "Usage: make engine HUGGINGFACE_TOKEN=hf_xxx"; \
		exit 1; \
	fi
	$(CONTAINER_RUNTIME) run --rm \
		--device nvidia.com/gpu=all \
		--security-opt label=disable \
		-v yodai_model-data:/data/models \
		-v $(CURDIR)/deploy/scripts/build-engine.sh:/data/build-engine.sh:ro,Z \
		-e HUGGINGFACE_TOKEN=$(HUGGINGFACE_TOKEN) \
		dustynv/tensorrt_llm:0.12-r36.4.0 \
		bash /data/build-engine.sh

up:
	$(CONTAINER_RUNTIME) compose -f deploy/docker-compose.yml up -d

down:
	$(CONTAINER_RUNTIME) compose -f deploy/docker-compose.yml down

clean:
	rm -rf bin/
