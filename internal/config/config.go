package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ListenAddr     string  `yaml:"listen_addr" envconfig:"YODAI_LISTEN_ADDR"`
	InferenceURL   string  `yaml:"inference_url" envconfig:"YODAI_INFERENCE_URL"`
	InferenceModel string  `yaml:"inference_model" envconfig:"YODAI_INFERENCE_MODEL"`
	MaxTokens      int     `yaml:"max_tokens" envconfig:"YODAI_MAX_TOKENS"`
	Temperature    float32 `yaml:"temperature" envconfig:"YODAI_TEMPERATURE"`
	TopP           float32 `yaml:"top_p" envconfig:"YODAI_TOP_P"`
}

func defaults() Config {
	return Config{
		ListenAddr:     ":8080",
		InferenceURL:   "http://localhost:8000",
		InferenceModel: "tensorrt_llm",
		MaxTokens:      512,
		Temperature:    0.7,
		TopP:           0.9,
	}
}

func Load(path string) (*Config, error) {
	cfg := defaults()

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("parsing config file: %w", err)
		}
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("processing env config: %w", err)
	}

	return &cfg, nil
}
