package config

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

type Config struct {
	System System     `yaml:"system`
	Logger zap.Config `yaml:"logger"`
}

type System struct {
	ProfileDir string `yaml:"profileDir"`
}

func NewConfig(paths ...string) Config {
	cfg := Config{
		System: System{
			ProfileDir: "/opt/cure",
		},
		Logger: zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
			Encoding:         "console",
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "T",
				LevelKey:       "L",
				NameKey:        "N",
				CallerKey:      "C",
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "M",
				StacktraceKey:  "S",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalColorLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
		},
	}

	for _, path := range paths {
		file, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = yaml.Unmarshal(file, &cfg)
		if err != nil {
			fmt.Printf("[Warning] Failed to read config at path %s", path)
			continue
		}
		return cfg
	}
	return cfg
}
