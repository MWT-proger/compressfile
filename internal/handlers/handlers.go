package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

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

var AllovedExtension = map[string]int{
	"png":  1,
	"jpg":  1,
	"jpeg": 1,
}

// getWidthXHeight() Возвращает строку вида WidthXHeight (200x200)
func (d imageTransformationData) getWidthXHeight() string {
	return fmt.Sprintf("%sx%s", d.width, d.height)
}

// checkFileFormat() Проверяет формат файла указанного в pathFile
func (d imageTransformationData) checkFileFormat() error {
	var (
		splitPathFile = strings.Split(d.pathFile, ".")
		extension     = splitPathFile[len(splitPathFile)-1]
	)

	if _, ok := AllovedExtension[extension]; !ok {
		return errors.New("формат не разрешен")
	}
	return nil
}

// setFromQuery(query url.Values) Устанавливает переменны согласно полученному Query
func (d *imageTransformationData) setFromQuery(query url.Values) error {
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
	td.setFromQuery(req.URL.Query())

	if td.pathFile == "" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	if td.width == "" && td.height == "" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	err := td.checkFileFormat()
	if err != nil {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	imgStockByte, err := h.storage.Get(config.BacketNameStorage, td.pathFile)

	if err != nil {
		log.Println(err)
		log.Println("Стоковое изображение не найдено")

		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	log.Println("Стоковое изображение получено")

	opt := transform.ParseOptions(td.getWidthXHeight())

	pathFileNewImg := basePath + td.getWidthXHeight() + "/" + td.pathFile

	// TODO: Проверка на существует ли такого формата картинка
	imgStockByteNewFormat, _ := h.storage.Get(config.BacketNameStorage, pathFileNewImg)

	if imgStockByteNewFormat != nil {
		log.Println("Уже существует")
		res.Write(imgStockByteNewFormat)
		return
	}

	img, _ := transform.Transform(imgStockByte, opt)

	err = h.storage.Put(img, config.BacketNameStorage, pathFileNewImg)

	if err != nil {
		log.Println(err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	res.Write(img)

}
