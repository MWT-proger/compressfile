package configs

import "os"

type Config struct {
	HostServer           string `env:"SERVER_ADDRESS"`
	BacketNameStorage    string `env:"BACKET_NAME_S3_STORAGE"`
	EndpointURLS3Storage string `env:"AWS_S3_ENDPOINT_URL"`
}

var newConfig Config

// InitConfig() Присваивает локальной не импортируемой переменной newConfig базовые значения
// Вызывается один раз при старте проекта
func InitConfig() *Config {
	newConfig = Config{
		HostServer:           ":8080",
		BacketNameStorage:    "",
		EndpointURLS3Storage: "",
	}
	return &newConfig
}

// GetConfig() выводит не импортируемую переменную newConfig
func GetConfig() Config {
	return newConfig
}

// SetConfigFromEnv() Прсваевает полям значения из ENV
// Вызывается один раз при старте проекта
func SetConfigFromEnv() {
	if envBaseURLShortener := os.Getenv("SERVER_ADDRESS"); envBaseURLShortener != "" {
		newConfig.HostServer = envBaseURLShortener
	}
	if envBacketNameStorage := os.Getenv("BACKET_NAME_S3_STORAGE"); envBacketNameStorage != "" {
		newConfig.BacketNameStorage = envBacketNameStorage
	}
	if envEndpointURLS3Storage := os.Getenv("AWS_S3_ENDPOINT_URL"); envEndpointURLS3Storage != "" {
		newConfig.EndpointURLS3Storage = envEndpointURLS3Storage
	}
}
