package main

import (
	"errors"
	"log"

	"github.com/MWT-proger/compressfile/configs"
	"github.com/MWT-proger/compressfile/internal/handlers"
	"github.com/MWT-proger/compressfile/internal/router"
	"github.com/MWT-proger/compressfile/internal/s3storage"
	"github.com/MWT-proger/compressfile/internal/server"
)

// init() Инициализирует настройки проекта
func init() {
	configInit := configs.InitConfig()
	parseFlags(configInit)
	configs.SetConfigFromEnv()
}

// validateConfig() Проверяет обязательные параметры для старта проекта
// при не соответствии требованиям возвращает ошибку
func validateConfig() error {
	ErrBacketNameStorageNotFound := errors.New("необходимо указать имя корзины s3 хранилища")
	ErrEndpointURLS3StorageNotFound := errors.New("необходимо указать URL адрес s3 хранилища")

	conf := configs.GetConfig()
	if conf.BacketNameStorage == "" {
		err := ErrBacketNameStorageNotFound
		return err
	}

	if conf.EndpointURLS3Storage == "" {
		err := ErrEndpointURLS3StorageNotFound
		return err
	}
	return nil
}

// main() Основной файл проекта
func main() {

	if err := validateConfig(); err != nil {
		log.Fatalln(err)
	} else {
		var (
			s    = &s3storage.Storage{}
			h, _ = handlers.NewAPIHandler(s)
			r    = router.Router(h)
			serv = server.Server{Router: r}
		)

		if err := s.InitClientS3(); err != nil {
			panic(err)
		}

		if err := serv.Run(); err != nil {
			panic(err)
		}
	}

}
