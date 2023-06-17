package main

import (
	"flag"

	"github.com/MWT-proger/compressfile/configs"
)

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags(conf *configs.Config) {

	flag.StringVar(&conf.HostServer, "a", conf.HostServer, "Адрес и порт запуска сервера.")

	flag.StringVar(&conf.BucketNameStorage, "bucket", conf.BucketNameStorage, "Имя корзины в которой лежат файлы.")

	flag.Parse()
}
