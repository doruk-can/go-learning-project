package config

import (
	"fmt"
	"github.com/codingconcepts/env"
	"time"
)

// EnvVars represents environment variables
type EnvVars struct {
	Service     Service
	Mongo       Mongo
	HTTPServer  HTTPServer
	ExternalAPI ExternalAPIConfig
}

// Service represents service configurations
type Service struct {
	Name        string `env:"SERVICE_NAME" default:"foover"` // for demo purposes, otherwise required:"true"
	LogLevel    string `env:"SERVICE_LOG_LEVEL" default:"INFO"`
	Environment string `env:"SERVICE_ENVIRONMENT" default:"local"`
}

// Mongo represents mongo configurations
type Mongo struct {
	URI               string        `env:"MONGO_URI" default:"mongodb://localhost:27017"` // for demo purposes, otherwise required:"true"
	Database          string        `env:"MONGO_DATABASE" default:"foover-db"`            // for demo purposes, otherwise required:"true"
	ConnectTimeout    time.Duration `env:"MONGO_CONNECT_TIMEOUT" default:"10s"`
	MinPoolSize       uint64        `env:"MONGO_MIN_POOL_SIZE" default:"4"`
	MaxPoolSize       uint64        `env:"MONGO_MAX_POOL_SIZE" default:"100"`
	PingTimeout       time.Duration `env:"MONGO_PING_TIMEOUT" default:"10s"`
	ReadTimeout       time.Duration `env:"MONGO_READ_TIMEOUT" default:"10s"`
	WriteTimeout      time.Duration `env:"MONGO_WRITE_TIMEOUT" default:"5s"`
	DisconnectTimeout time.Duration `env:"MONGO_DISCONNECT_TIMEOUT" default:"5s"`
}

// HTTPServer represents http server configurations
type HTTPServer struct {
	Address         string        `env:"HTTP_SERVER_ADDRESS" default:":8080"`
	ReadTimeout     time.Duration `env:"HTTP_SERVER_READ_TIMEOUT" default:"15s"`
	WriteTimeout    time.Duration `env:"HTTP_SERVER_WRITE_TIMEOUT" default:"30s"`
	IdleTimeout     time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT" default:"15s"`
	MaxHeaderBytes  int           `env:"HTTP_SERVER_MAX_HEADER_BYTES" default:"1048576"`
	ShutdownTimeout time.Duration `env:"HTTP_SERVER_SHUTDOWN_TIMEOUT" default:"10s"`
}

// ExternalAPIConfig represents external api configurations
type ExternalAPIConfig struct {
	ProductAPIURL string `env:"EXTERNAL_API_PRODUCT_URL" default:"https://amperoid.tenants.foodji.io/machines/4bf115ee-303a-4089-a3ea-f6e7aae0ab94"` // for demo purposes, otherwise required:"true"
}

// LoadEnvVars loads and returns environment variables
func LoadEnvVars() (*EnvVars, error) {
	s := Service{}
	if err := env.Set(&s); err != nil {
		return nil, fmt.Errorf("loading service environment variables failed, %s", err.Error())
	}

	m := Mongo{}
	if err := env.Set(&m); err != nil {
		return nil, fmt.Errorf("loading mongo environment variables failed, %s", err.Error())
	}

	hs := HTTPServer{}
	if err := env.Set(&hs); err != nil {
		return nil, fmt.Errorf("loading http server environment variables failed, %s", err.Error())
	}

	ea := ExternalAPIConfig{}
	if err := env.Set(&ea); err != nil {
		return nil, fmt.Errorf("loading external api environment variables failed, %s", err.Error())
	}

	ev := &EnvVars{
		Service:     s,
		Mongo:       m,
		HTTPServer:  hs,
		ExternalAPI: ea,
	}

	return ev, nil
}
