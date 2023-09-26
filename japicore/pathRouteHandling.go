package japicore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/JackalLabs/jackalapi/jutils"
	"github.com/JackalLabs/jackalgo/handlers/file_io_handler"
	"github.com/uptrace/bunrouter"
)

func downloadByPathCore(fileIo *file_io_handler.FileIoHandler, operatingRoot string) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		location := req.Param("location")
		if len(location) == 0 {
			warning := "Failed to get Location"
			return jutils.ProcessCustomHttpError("DownloadByPathHandler", warning, 404, w)
		}

		uniquePath := readUniquePath(req)
		if len(uniquePath) > 0 {
			operatingRoot += "/" + uniquePath
		}
		operatingRoot += "/" + location

		handler, err := fileIo.DownloadFile(operatingRoot)
		if err != nil {
			return err
		}

		fileBytes := handler.GetFile().Buffer().Bytes()
		_, err = w.Write(fileBytes)
		if err != nil {
			jutils.ProcessError("WWriteError for DownloadByPathHandler", err)
		}
		return nil
	}
}

func deleteByPathCore(fileIo *file_io_handler.FileIoHandler, operatingRoot string) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		filename := req.Param("filename")
		if len(filename) == 0 {
			warning := "Failed to get FileName"
			return jutils.ProcessCustomHttpError("processUpload", warning, 400, w)
		}

		location := req.Param("location")
		if len(location) == 0 {
			warning := "Failed to get Location"
			return jutils.ProcessCustomHttpError("DownloadByPathHandler", warning, 404, w)
		}

		cleanFilename := strings.ReplaceAll(filename, "/", "_")
		fmt.Println(cleanFilename)

		uniquePath := readUniquePath(req)
		if len(uniquePath) > 0 {
			operatingRoot += "/" + uniquePath
		}
		operatingRoot += "/" + location

		folder, err := fileIo.DownloadFolder(operatingRoot)
		if err != nil {
			jutils.ProcessHttpError("DeleteFile", err, 404, w)
			return err
		}

		err = fileIo.DeleteTargets([]string{cleanFilename}, folder)
		if err != nil {
			jutils.ProcessHttpError("DeleteFile", err, 500, w)
			return err
		}

		message := createJsonResponse("Deletion complete")
		condensedWriteJSON(w, message)
		return nil
	}
}

func ImportHandler(fileIo *file_io_handler.FileIoHandler, queue *ScrapeQueue) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_BULK_ROOT", "s/JAPI/Bulk")

		uniquePath := readUniquePath(req)
		if len(uniquePath) > 0 {
			operatingRoot += "/" + uniquePath
		}

		var data fileScrape
		source := req.Header.Get("J-Source-Path")
		operatingRoot += "/" + source

		err := json.NewDecoder(req.Body).Decode(&data)
		if err != nil {
			jutils.ProcessHttpError("JSONDecoder", err, 500, w)
			return err
		}

		var wg sync.WaitGroup

		for _, target := range data.Targets {
			wg.Add(1)
			queue.Push(fileIo, w, &wg, operatingRoot, target, source)
		}

		wg.Wait()

		message := createJsonResponse("Import complete")
		condensedWriteJSON(w, message)
		return nil
	}
}

func DownloadFromBulkByPathHandler(fileIo *file_io_handler.FileIoHandler) bunrouter.HandlerFunc {
	operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_BULK_ROOT", "s/JAPI/Bulk")
	return downloadByPathCore(fileIo, operatingRoot)
}

func DownloadByPathHandler(fileIo *file_io_handler.FileIoHandler) bunrouter.HandlerFunc {
	operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_OP_ROOT", "s/JAPI")
	return downloadByPathCore(fileIo, operatingRoot)
}

func UploadByPathHandler(fileIo *file_io_handler.FileIoHandler, queue *FileIoQueue) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_OP_ROOT", "s/JAPI")
		var byteBuffer bytes.Buffer
		var wg sync.WaitGroup
		wg.Add(1)
		WorkingFileSize := 32 << 30

		envSize := jutils.LoadEnvVarOrFallback("JAPI_MAX_FILE", "")
		if len(envSize) > 0 {
			envParse, err := strconv.Atoi(envSize)
			if err != nil {
				return err
			}
			WorkingFileSize = envParse
		}
		MaxFileSize := int64(WorkingFileSize)

		// ParseMultipartForm parses a request body as multipart/form-data
		err := req.ParseMultipartForm(MaxFileSize) // MAX file size lives here
		if err != nil {
			jutils.ProcessHttpError("ParseMultipartForm", err, 400, w)
			return err
		}

		// Retrieve the file from form data
		file, head, err := req.FormFile("file")
		if err != nil {
			jutils.ProcessHttpError("FormFileFile", err, 400, w)
			return err
		}

		uniquePath := readUniquePath(req)
		if len(uniquePath) > 0 {
			operatingRoot += "/" + uniquePath
		}

		subFolder := req.FormValue("subfolder")
		if len(subFolder) > 0 {
			operatingRoot += "/" + subFolder
		}

		_, err = io.Copy(&byteBuffer, file)
		if err != nil {
			jutils.ProcessHttpError("MakeByteBuffer", err, 500, w)
			return err
		}

		fid := processUpload(w, fileIo, byteBuffer.Bytes(), head.Filename, operatingRoot, queue)
		if len(fid) == 0 {
			warning := "Failed to get FID"
			return jutils.ProcessCustomHttpError("processUpload", warning, 500, w)
		}

		successfulUpload := UploadResponse{
			FID: fid,
		}
		err = json.NewEncoder(w).Encode(successfulUpload)
		if err != nil {
			jutils.ProcessHttpError("JSONSuccessEncode", err, 500, w)
			return err
		}

		message := createJsonResponse("Upload complete")
		condensedWriteJSON(w, message)
		return nil
	}
}

func DeleteFromBulkByPathHandler(fileIo *file_io_handler.FileIoHandler) bunrouter.HandlerFunc {
	operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_BULK_ROOT", "s/JAPI/Bulk")
	return deleteByPathCore(fileIo, operatingRoot)
}

func DeleteByPathHandler(fileIo *file_io_handler.FileIoHandler) bunrouter.HandlerFunc {
	operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_OP_ROOT", "s/JAPI")
	return deleteByPathCore(fileIo, operatingRoot)
}
