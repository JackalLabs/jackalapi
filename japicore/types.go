package japicore

import (
	"net/http"

	"github.com/JackalLabs/jackalgo/handlers/wallet_handler"

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
	Fids    []string `json:"fids"`
	Message string   `json:"message"`
	Module  string   `json:"module"`
	Version string   `json:"version"`
}

func createJsonResponse(message string, fids []string) jsonResponse {
	return jsonResponse{
		Fids:    fids,
		Message: message,
		Module:  Module,
		Version: Version,
	}
}

type JApiCore struct {
	FileIo      *file_io_handler.FileIoHandler
	FileIoQueue *FileIoQueue
	ScrapeQueue *ScrapeQueue
	Wallet      *wallet_handler.WalletHandler
}

func InitJApiCore() *JApiCore {
	wallet, fileIo := InitWalletSession()
	fileIoQueue := NewFileIoQueue()
	scrapeQueue := NewScrapeQueue(fileIoQueue)

	go fileIoQueue.Listen()

	core := JApiCore{
		FileIo:      fileIo,
		FileIoQueue: fileIoQueue,
		ScrapeQueue: scrapeQueue,
		Wallet:      wallet,
	}

	return &core
}
