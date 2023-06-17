package configs

import "os"

var AllovedExtension = map[string]int{
	"png":  1,
	"jpg":  1,
	"jpeg": 1,
}

type Config struct {
	HostServer           string `env:"SERVER_ADDRESS"`
	BucketNameStorage    string `env:"AWS_S3_BACKET_NAME"`
	EndpointURLS3Storage string `env:"AWS_S3_ENDPOINT_URL"`
	AllovedExtension     map[string]int
}

var newConfig Config

// InitConfig() Присваивает локальной не импортируемой переменной newConfig базовые значения
// Вызывается один раз при старте проекта
func InitConfig() *Config {
	newConfig = Config{
		HostServer:           ":8080",
		BucketNameStorage:    "",
		EndpointURLS3Storage: "",
		AllovedExtension:     AllovedExtension,
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
	if envServerAddres := os.Getenv("SERVER_ADDRESS"); envServerAddres != "" {
		newConfig.HostServer = envServerAddres
	}
	if envBucketNameStorage := os.Getenv("AWS_S3_BUCKET_NAME"); envBucketNameStorage != "" {
		newConfig.BucketNameStorage = envBucketNameStorage
	}
	if envEndpointURLS3Storage := os.Getenv("AWS_S3_ENDPOINT_URL"); envEndpointURLS3Storage != "" {
		newConfig.EndpointURLS3Storage = envEndpointURLS3Storage
	}
}
