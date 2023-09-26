package japicore

import (
	"net/http"

	"github.com/uptrace/bunrouter"

	"github.com/JackalLabs/jackalgo/handlers/file_io_handler"
)

type Handlers map[string]*func(w http.ResponseWriter, r bunrouter.Request, queue *FileIoQueue, fileIo *file_io_handler.FileIoHandler) error

type UploadResponse struct {
	FID string `json:"fid"`
}

type fileScrape struct {
	Targets []string `json:"targets"`
}

type jsonResponse struct {
	Module  string `json:"module"`
	Version string `json:"version"`
	Message string `json:"message"`
}

func createJsonResponse(message string) jsonResponse {
	return jsonResponse{
		Message: message,
		Module:  Module,
		Version: Version,
	}
}
