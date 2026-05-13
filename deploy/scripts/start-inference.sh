#!/bin/bash
set -e

echo "Starting TensorRT-LLM OpenAI server..."
echo "Engine: ${ENGINE_DIR}"
echo "Tokenizer: ${TOKENIZER_DIR}"
echo "Listening on ${HOST}:${PORT}"

python3 /opt/TensorRT-LLM/examples/apps/openai_server.py \
    "${ENGINE_DIR}" \
    --tokenizer "${TOKENIZER_DIR}" \
    --host "${HOST}" \
    --port "${PORT}"
