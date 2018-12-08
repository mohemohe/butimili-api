package configs

import (
	"github.com/vrischmann/envconfig"
)

type (
	Env struct {
		Echo struct {
			Env     string `envconfig:"default=debug"`
			Address string `envconfig:"default=0.0.0.0:1323"`
		}
		Mongo struct {
			Address  string `envconfig:"default=mongodb://localhost:27017"`
			Database string `envconfig:"default=butimili"`
			SSL      bool   `envconfig:"default=false"`
		}
		Sign struct {
			Secret string `envconfig:"default=changeme"`
		}
		Encrypt struct {
			Secret string `envconfig:"default=changeme"`
		}
	}
)

var envCached bool
var env Env

func GetEnv() Env {
	if !envCached {
		if err := envconfig.Init(&env); err != nil {
			panic(err)
		}
		envCached = true
	}
	return env
}
