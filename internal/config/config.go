// Package config defines the cure options
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

// Config is the main config object
type Config struct {
	System System     `yaml:"system"`
	Logger zap.Config `yaml:"logger"`
}

// System contains info about system wide profile
type System struct {
	ProfileDir string `yaml:"profileDir"`
}

// NewConfig creates a new Config and load with given paths
func NewConfig(paths ...string) Config {
	cfg := Config{
		System: System{
			ProfileDir: "/opt/cure",
		},
		Logger: zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.WarnLevel),
			Encoding:         "console",
			OutputPaths:      []string{"stderr", "/tmp/app.log"},
			ErrorOutputPaths: []string{"stderr", "/tmp/app-error.log"},
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
		file, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
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
