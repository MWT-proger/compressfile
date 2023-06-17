package errors

import (
	"github.com/MWT-proger/compressfile/configs"
)

type ErrorBacketNameS3StorageNotFound struct{}
type ErrorEndpointURLS3StorageNotFound struct{}
type ErrorNoSuchKeyInS3Storage struct{}
type ErrorFileExtensionNotAllowed struct{}

func (m *ErrorBacketNameS3StorageNotFound) Error() string {
	return "необходимо указать имя корзины s3 хранилища"
}

func (m *ErrorEndpointURLS3StorageNotFound) Error() string {
	return "необходимо указать URL адрес s3 хранилища"
}

func (m *ErrorNoSuchKeyInS3Storage) Error() string {
	return "данный ключ не существует в s3 storage"
}

func (m *ErrorFileExtensionNotAllowed) Error() string {
	allovedExtensionString := ""
	for k := range configs.GetConfig().AllovedExtension {
		allovedExtensionString += k + ","
	}
	return "расширение файла не разрешено. допустимые расширения: " + allovedExtensionString
}
