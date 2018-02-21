package config

import (
	"gopkg.in/tokopedia/logging.v1/tracer"
)

var Config = tracer.Config{
	Port:    8700,
	Enabled: true,
	TTL:     3600,
	Name:    "jaeger",
}
