package config

import "time"

func LoadDefConf() *Config {
	s := GrpcServer{55000, time.Second * 10}

	d := Config{
		Env:         "local",
		StoragePath: "./internal/storages/exchanger.db",
		GrpcServ:    s,
	}

	return &d
}
