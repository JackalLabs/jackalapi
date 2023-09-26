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

func ImportHandler(fileIo *file_io_handler.FileIoHandler, queue *ScrapeQueue) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_BULK_ROOT", "s/JAPI/Bulk")

		var data fileScrape
		source := req.Header.Get("J-Source-Path")

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

func DownloadByPathHandler(fileIo *file_io_handler.FileIoHandler) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		id := req.Param("id")
		if len(id) == 0 {
			warning := "Failed to get FileName"
			return jutils.ProcessCustomHttpError("processUpload", warning, 404, w)
		}
		fid := strings.ReplaceAll(id, "/", "_")

		handler, err := fileIo.DownloadFileFromFid(fid)
		if err != nil {
			return err
		}

		fileBytes := handler.GetFile().Buffer().Bytes()
		_, err = w.Write(fileBytes)
		if err != nil {
			jutils.ProcessError("WWriteError for DownloadByFidHandler", err)
		}
		return nil
	}
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

func DeleteByPathHandler(fileIo *file_io_handler.FileIoHandler, queue *FileIoQueue) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		id := req.Param("id")
		if len(id) == 0 {
			warning := "Failed to get FileName"
			return jutils.ProcessCustomHttpError("processUpload", warning, 400, w)
		}

		fid := strings.ReplaceAll(id, "/", "_")
		fmt.Println(fid)

		folder, err := fileIo.DownloadFolder(queue.GetRoot("bulk"))
		if err != nil {
			jutils.ProcessHttpError("DeleteFile", err, 404, w)
			return err
		}

		err = fileIo.DeleteTargets([]string{fid}, folder)
		if err != nil {
			jutils.ProcessHttpError("DeleteFile", err, 500, w)
			return err
		}

		message := createJsonResponse("Deletion complete")
		condensedWriteJSON(w, message)
		return nil
	}
}
