package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	GrpcServ    GrpcServer `yaml:"grpc_server"`
}

type GrpcServer struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoadConf() *Config {
	var cfg Config

	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
