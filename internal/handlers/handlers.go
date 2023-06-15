package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

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

type imageTransformationData struct {
	pathFile string
	width    string
	height   string
}

// getWidthXHeight() Возвращает строку вида WidthXHeight (200x200)
func (d imageTransformationData) getWidthXHeight() string {
	return fmt.Sprintf("%sx%s", d.width, d.height)
}

// checkFileFormat() Проверяет формат файла указанного в pathFile
func (d imageTransformationData) checkFileFormat() error {
	// TODO: Проверка на определенные форматы картинок
	return nil
}

// getFromQuery(query url.Values) Устанавливает переменны согласно полученному Query
func (d *imageTransformationData) getFromQuery(query url.Values) error {
	d.pathFile = query.Get("pathFile")
	d.width = query.Get("width")
	d.height = query.Get("height")

	return nil
}

// TransformImage() Плучает картинку по указанному пути
// Пример "collections/c547ffa8-9d26-4141-a54d-f2f4ae4d8153/28d1d5ec-ddc2-4edb-8adf-37e75031e109/tokens/cfad535a28a84b198811d225a61f566a.png"
// трансформирует изображение
// и сохраняет обратно в хранилище добавляя в путь файла параметры трансформации
func (h *APIHandler) TransformImage(res http.ResponseWriter, req *http.Request) {

	var (
		config   = configs.GetConfig()
		basePath = "GoCompress/v1/"
		td       = imageTransformationData{}
	)
	td.getFromQuery(req.URL.Query())

	if td.pathFile == "" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	if td.width == "" && td.height == "" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	td.checkFileFormat()

	imgStockByte, err := h.storage.Get(config.BacketNameStorage, td.pathFile)

	if err != nil {
		log.Println(err)
		log.Println("Стоковое изображение не найдено")

		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	log.Println("Стоковое изображение получено")

	opt := transform.ParseOptions(td.getWidthXHeight())

	// TODO: Проверка на существует ли такого формата картинка
	imgStockByteNewFormat, err := h.storage.Get(config.BacketNameStorage, basePath+td.pathFile)

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

	err = h.storage.Put(img, config.BacketNameStorage, basePath+td.pathFile)

	if err != nil {
		log.Println(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	res.Write(img)

}
