package config

import "time"

func LoadDefConf() *Config {
	s := GrpcServer{
		Port: 55000,
		MaxConnectionIdle: time.Second * 15,
		MaxConnectionAge: time.Second * 60,
		MaxConnectionAgeGrace: time.Second * 10,
		Time: time.Second * 15,
		Timeout: time.Second * 5,               
	}

	d := Config{
		Env:         "local",
		StoragePath: "./internal/storages/exchanger.db",
		GrpcServ:    s,
	}

	return &d
}
