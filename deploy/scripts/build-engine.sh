#!/bin/bash
set -e

MODEL_NAME="${MODEL_NAME:-meta-llama/Meta-Llama-3.1-8B-Instruct}"
MODEL_DIR="${MODEL_DIR:-/data/models/hf}"
ENGINE_DIR="${ENGINE_DIR:-/data/models/engine}"
TOKENIZER_DIR="${TOKENIZER_DIR:-/data/models/tokenizer}"
QUANT_TYPE="${QUANT_TYPE:-int4_awq}"
MAX_SEQ_LEN="${MAX_SEQ_LEN:-2048}"
MAX_BATCH_SIZE="${MAX_BATCH_SIZE:-1}"

if [ -z "${HUGGINGFACE_TOKEN}" ]; then
    echo "Error: HUGGINGFACE_TOKEN is required"
    echo "Usage: HUGGINGFACE_TOKEN=hf_xxx bash build-engine.sh"
    exit 1
fi

echo "=== YodAI Engine Builder ==="
echo "Model: ${MODEL_NAME}"
echo "Quantization: ${QUANT_TYPE}"
echo "Max sequence length: ${MAX_SEQ_LEN}"
echo ""

# Step 1: Download model
echo "[1/3] Downloading model from Hugging Face..."
huggingface-cli login --token "${HUGGINGFACE_TOKEN}"
huggingface-cli download "${MODEL_NAME}" --local-dir "${MODEL_DIR}"

# Copy tokenizer files separately for the server
echo "[2/3] Preparing tokenizer..."
mkdir -p "${TOKENIZER_DIR}"
cp "${MODEL_DIR}"/tokenizer* "${TOKENIZER_DIR}/" 2>/dev/null || true
cp "${MODEL_DIR}"/special_tokens_map.json "${TOKENIZER_DIR}/" 2>/dev/null || true

# Step 3: Quantize and build engine
echo "[3/3] Building TensorRT-LLM engine (this may take 15-30 minutes)..."
mkdir -p "${ENGINE_DIR}"

python3 /opt/TensorRT-LLM/examples/quantization/quantize.py \
    --model_dir "${MODEL_DIR}" \
    --output_dir "${ENGINE_DIR}/quantized" \
    --dtype float16 \
    --qformat "${QUANT_TYPE}" \
    --calib_size 32

trtllm-build \
    --checkpoint_dir "${ENGINE_DIR}/quantized" \
    --output_dir "${ENGINE_DIR}" \
    --gemm_plugin float16 \
    --max_batch_size "${MAX_BATCH_SIZE}" \
    --max_input_len "${MAX_SEQ_LEN}" \
    --max_seq_len "${MAX_SEQ_LEN}"

echo ""
echo "=== Engine build complete ==="
echo "Engine files: ${ENGINE_DIR}"
echo "Tokenizer files: ${TOKENIZER_DIR}"
echo "Ready to start the inference server."
