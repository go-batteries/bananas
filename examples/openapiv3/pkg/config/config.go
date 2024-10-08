package config

import (
	"os"

	"github.com/go-batteries/diaper"
	"github.com/rs/zerolog/log"
)

type AppConfig struct {
	ServerID    string
	Environment string
	Dsn         string
	DbName      string
	RedisURL    string
	AppPort     string
	GrpcPort    string
}

func Load(envFile string) AppConfig {
	providers := diaper.BuildProviders(diaper.EnvProvider{})
	loader := diaper.DiaperConfig{
		DefaultEnvFile: "app.env",
		Providers:      providers,
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "dev"
	}

	cfgMap, err := loader.ReadFromFile(env, envFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config from " + envFile)
	}

	return AppConfig{
		Environment: env,
		ServerID:    cfgMap.MustGet("server_id").(string),
		Dsn:         cfgMap.MustGet("dsn").(string),
		DbName:      cfgMap.MustGet("db_name").(string),
		RedisURL:    cfgMap.MustGet("redis_url").(string),
		AppPort:     cfgMap.MustGet("app_port").(string),
		GrpcPort:    cfgMap.MustGet("grpc_port").(string),
	}
}
