package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MWT-proger/compressfile/configs"
	"github.com/MWT-proger/compressfile/internal/s3storage"
	"github.com/MWT-proger/compressfile/internal/transform"
)

type APIHandler struct {
	storage s3storage.OperationStorager
}

func NewAPIHandler(s s3storage.OperationStorager) (h *APIHandler, err error) {
	return &APIHandler{s}, err
}

// TransformImage() Плучает картинку по указанному пути
// Пример "collections/c547ffa8-9d26-4141-a54d-f2f4ae4d8153/28d1d5ec-ddc2-4edb-8adf-37e75031e109/tokens/cfad535a28a84b198811d225a61f566a.png"
// трансформирует изображение
// и сохраняет обратно в хранилище добавляя в путь файла параметры трансформации
func (h *APIHandler) TransformImage(res http.ResponseWriter, req *http.Request) {
	config := configs.GetConfig()

	var (
		pathFile = req.URL.Query().Get("pathFile")
		width    = req.URL.Query().Get("width")
		height   = req.URL.Query().Get("height")
		basePath = "GoCompress/v1/"
	)

	if pathFile == "" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	if width == "" && height == "" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	// TODO: Проверка на определенные форматы картинок

	imgStockByte, err := h.storage.Get(config.BacketNameStorage, pathFile)

	if err != nil {
		log.Println(err)
		log.Println("Стоковое изображение не найдено")

		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	log.Println("Стоковое изображение получено")

	resizeParams := fmt.Sprintf("%sx%s", width, height)
	opt := transform.ParseOptions(resizeParams)

	// TODO: Проверка на существует ли такого формата картинка
	imgStockByteNewFormat, err := h.storage.Get(config.BacketNameStorage, basePath+pathFile)

	if err != nil {
		log.Println(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	if imgStockByteNewFormat != nil {
		log.Println("Уже существует")
		res.Write(imgStockByteNewFormat)
		return
	}

	img, _ := transform.Transform(imgStockByte, opt)

	err = h.storage.Put(img, config.BacketNameStorage, basePath+pathFile)

	if err != nil {
		log.Println(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	res.Write(img)

}
