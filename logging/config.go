package logging

import "github.com/dubbikins/envy"

var logConfig = &LogConfig{}

func WithLoggerConfig(options ...func(config *LogConfig)) {
	for _, f := range options {
		f(logConfig)
	}
}

func FromEnv(config *LogConfig) {
	envy.Unmarshal(config)
}

type LogConfig struct {
	WithColor bool `env:"GLAM_LOG_WITH_COLOR"`
}
