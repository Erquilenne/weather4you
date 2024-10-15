package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database        DatabaseConfig `mapstructure:"database"`
	Logger          Logger         `mapstructure:"logger"`
	Server          ServerConfig   `mapstructure:"server"`
	StartCities     []string       `mapstructure:"start_cities"`
	WeatherApiToken string         `mapstructure:"weather_api_token"`
	Jaeger          Jaeger         `mapstructure:"jaeger"`
	Metrics         Metrics        `mapstructure:"metrics"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type ServerConfig struct {
	Mode         string        `mapstructure:"mode"`
	AppVersion   string        `mapstructure:"app_version"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	Port         string        `mapstructure:"port"`
	PprofPort    string        `mapstructure:"pprof_port"`
	Debug        bool          `mapstructure:"debug"`
}

type Logger struct {
	Development       bool   `mapstructure:"development"`
	DisableCaller     bool   `mapstructure:"disable_caller"`
	DisableStacktrace bool   `mapstructure:"disable_stacktrace"`
	Encoding          string `mapstructure:"encoding"`
	Level             string `mapstructure:"level"`
}

type Jaeger struct {
	Host        string `mapstructure:"host"`
	ServiceName string `mapstructure:"service_name"`
	LogSpans    bool   `mapstructure:"log_spans"`
}

type Metrics struct {
	URL         string `mapstructure:"url"`
	ServiceName string `mapstructure:"service_name"`
}

func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.SetConfigType("json")
	v.AddConfigPath(".")

	v.AddConfigPath("./config")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, err
			// return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var cfg Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &cfg, nil
}
