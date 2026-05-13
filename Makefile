.PHONY: build run test clean engine docker-build up down dev tidy

BINARY := bin/yodai
GOARCH ?= arm64

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

docker-build:
	docker compose -f deploy/docker-compose.yml build

engine:
	@if [ -z "$(HUGGINGFACE_TOKEN)" ]; then \
		echo "Usage: make engine HUGGINGFACE_TOKEN=hf_xxx"; \
		exit 1; \
	fi
	docker run --rm --runtime nvidia \
		-v yodai_model-data:/data/models \
		-e HUGGINGFACE_TOKEN=$(HUGGINGFACE_TOKEN) \
		dustynv/tensorrt_llm:0.12-r36.4.0 \
		bash /data/build-engine.sh

up:
	docker compose -f deploy/docker-compose.yml up -d

down:
	docker compose -f deploy/docker-compose.yml down

clean:
	rm -rf bin/
