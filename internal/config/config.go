package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	GrpcServ    GrpcServer `yaml:"grpc_server"`
}

type GrpcServer struct {
	Port              int           `yaml:"port"`
	ConnectionTimeout time.Duration `yaml:"connectionTimeout"`
}

// you can pass a file to the configuration to run it or run it with default parameters
func MustLoadConf() *Config {
	var cfg Config
	var filePath string

	flag.StringVar(&filePath, "config", "", "path to config file")
	flag.Parse()
	
	switch filePath {
	case "":
		panic("no configuration is specified in the config flag")
	case "default":
		log.Println("ATTENTION!!! The server will be started with the default configuration")
		defCfg := LoadDefConf()
		return defCfg
	default:
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			panic("config file does not exist: " + filePath)
		}
	}

	err := cleanenv.ReadConfig(filePath, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
