package handlers

import (
	"log"
	"net/http"

	"github.com/MWT-proger/compressfile/configs"
	"github.com/MWT-proger/compressfile/internal/s3storage"
)

type APIHandler struct {
	storage s3storage.OperationStorager
}

func NewAPIHandler(s s3storage.OperationStorager) (h *APIHandler, err error) {
	return &APIHandler{s}, err
}

// TransformImage() Плучает картинку по указанному пути
// трансформирует изображение
// и сохраняет обратно в хранилище добавляя в путь файла параметры трансформации
func (h *APIHandler) TransformImage(res http.ResponseWriter, req *http.Request) {
	config := configs.GetConfig()
	testKey := "collections/c547ffa8-9d26-4141-a54d-f2f4ae4d8153/28d1d5ec-ddc2-4edb-8adf-37e75031e109/tokens/cfad535a28a84b198811d225a61f566a.png"
	log.Print(config.BacketNameStorage, testKey)
	file, err := h.storage.UploadFileToServer(config.BacketNameStorage, testKey)

	if err != nil {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	res.Write(file)

}