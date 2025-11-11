package configs

import (
	"L2_18/configs/loader"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type HttpConfig struct {
	Port         string        `validate:"required"`
	ReadTimeout  time.Duration `validate:"required"`
	WriteTimeout time.Duration `validate:"required"`
	IdleTimeout  time.Duration `validate:"required"`
}

type Config struct {
	HTTP HttpConfig
	Env  string
}

func MustLoad(loader loader.ConfigLoader) *Config {
	env := os.Getenv("APP_ENV")
	if env == "" {
		envFlag := flag.String("env", "dev", "Environment type")
		flag.Parse()
		env = *envFlag
	}

	const op = "configs.MustLoad"
	envs, err := loader.Load()
	if err != nil {
		log.Fatalf("%s: config load failed: %+v", op, err)
	}
	cfg := &Config{
		HTTP: HttpConfig{
			Port:         envs["HTTP_PORT"],
			ReadTimeout:  getEnvAsDuration(envs["HTTP_READ_TIMEOUT"], 10*time.Second),
			WriteTimeout: getEnvAsDuration(envs["HTTP_WRITE_TIMEOUT"], 10*time.Second),
			IdleTimeout:  getEnvAsDuration(envs["HTTP_WRITE_TIMEOUT"], 60*time.Second),
		},
		Env: env,
	}

	if err := validateConfig(cfg); err != nil {
		log.Fatalf("%s: error validation config: %+v", op, err)
	}

	return cfg
}

func validateConfig(cfg *Config) error {
	if cfg.HTTP.Port == "" || cfg.HTTP.ReadTimeout <= 0*time.Second || cfg.HTTP.WriteTimeout <= 0*time.Second ||
		cfg.HTTP.IdleTimeout <= 0*time.Second {
		return fmt.Errorf("incorrect http config fields")
	}
	return nil
}

func getEnvAsDuration(strValue string, defaultValue time.Duration) time.Duration {
	const op = "configs.getEnvAsDuration"
	if strValue == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(strValue)
	if err != nil {
		log.Printf("%s:forbidden value for %s, using default: %v", op,
			strValue, defaultValue)
		return defaultValue
	}
	return value
}

//
//func getEnvAsInt(strValue string, defaultValue int) int {
//	const op = "configs.getEnvAsInt"
//	if strValue == "" {
//		return defaultValue
//	}
//	value, err := strconv.Atoi(strValue)
//	if err != nil {
//		log.Printf("%s:forbidden value for %s, using default: %v", op, strValue,
//			defaultValue)
//		return defaultValue
//	}
//	return value
//}
//
//func getEnvAsBool(strValue string, defaultValue bool) bool {
//	const op = "configs.getEnvAsBool"
//	if strValue == "" {
//		return defaultValue
//	}
//	value, err := strconv.ParseBool(strValue)
//	if err != nil {
//		log.Printf("%s:forbidden value for %s, using default: %v", op, strValue, defaultValue)
//		return defaultValue
//	}
//	return value
//}
