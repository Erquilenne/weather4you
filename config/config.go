package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database        DatabaseConfig `json:"database"`
	Logger          Logger         `json:"logger"`
	Server          ServerConfig   `json:"server"`
	StartCities     []string       `json:"start_cities"`
	WeatherApiToken string         `json:"weather_api_token"`
	Jaeger          Jaeger         `json:"jaeger"`
	Metrics         Metrics        `json:"metrics"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type ServerConfig struct {
	Mode         string `json:"mode"`
	AppVersion   string `json:"app_version"`
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Port         string
	PprofPort    string
	Debug        bool
}

type Logger struct {
	Development       bool   `json:"development"`
	DisableCaller     bool   `json:"disable_caller"`
	DisableStacktrace bool   `json:"disable_stacktrace"`
	Encoding          string `json:"encoding"`
	Level             string `json:"level"`
}

type Jaeger struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

// Metrics config
type Metrics struct {
	URL         string
	ServiceName string
}

func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(filename) // Укажите имя файла конфигурации без расширения
	v.SetConfigType("json")   // Укажите тип файла конфигурации
	v.AddConfigPath(".")      // Укажите путь к конфигурационному файлу

	// Можно указать дополнительные пути
	v.AddConfigPath("./config")

	v.AutomaticEnv() // Поддержка переменных окружения

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
