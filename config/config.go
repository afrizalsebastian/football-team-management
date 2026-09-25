package config

import (
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/spf13/viper"
)

var (
	ErrOsGetPwd            = errors.New("failed to get working directory")
	ErrReadConfigFile      = errors.New("failed to read config file")
	ErrUnmarshalConfigFile = errors.New("failed to read config file")
)

type Config struct {
	ServerPort      int           `mapstructure:"SERVER_PORT"`
	DBUser          string        `mapstructure:"DB_USER"`
	DBPassword      string        `mapstructure:"DB_PASSWORD"`
	DBHost          string        `mapstructure:"DB_HOST"`
	DBPort          string        `mapstructure:"DB_PORT"`
	DBName          string        `mapstructure:"DB_NAME"`
	JwtSecret       string        `mapstructure:"JWT_SECRET"`
	JwtTTL          time.Duration `mapstructure:"JWT_TTL"`
	SuperadminToken []string      `mapstructure:"SUPERADMIN_TOKEN"`
}

func bindEnvs(v *viper.Viper, config interface{}) {
	cfgType := reflect.TypeOf(config).Elem()
	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)
		envKey := field.Tag.Get("mapstructure")
		if envKey != "" {
			_ = v.BindEnv(envKey)
		}
	}
}

func Init(wd string) (*Config, error) {
	var appConfig Config
	v := viper.New()

	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	v.AddConfigPath(wd)

	if err := v.ReadInConfig(); err != nil {
		var errConfigFileNotFound viper.ConfigFileNotFoundError
		if errors.As(err, &errConfigFileNotFound) {
			log.Println("No .env file found")
		} else {
			log.Println("Failed to read .env file")
			return nil, ErrReadConfigFile
		}
	}

	bindEnvs(v, &appConfig)

	if err := v.Unmarshal(&appConfig); err != nil {
		log.Println("error to unmarshal config file")
		return nil, ErrUnmarshalConfigFile
	}

	return &appConfig, nil
}
