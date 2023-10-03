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
	"github.com/uptrace/bunrouter"
)

func (j JApiCore) downloadByPathCore(operatingRoot string, reportFunc func(num int64)) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		location := req.Param("location")
		if len(location) == 0 {
			warning := "Failed to get Location"
			return jutils.ProcessCustomHttpError("BasicDownloadByPathHandler", warning, 404, w)
		}

		uniquePath := readUniquePath(req)
		if len(uniquePath) > 0 {
			operatingRoot += "/" + uniquePath
		}
		operatingRoot += "/" + location

		handler, err := j.FileIo.DownloadFile(operatingRoot)
		if err != nil {
			return err
		}

		size := handler.GetFile().Details.Size
		reportFunc(size)

		fileBytes := handler.GetFile().Buffer().Bytes()
		_, err = w.Write(fileBytes)
		if err != nil {
			jutils.ProcessError("WWriteError for BasicDownloadByPathHandler", err)
		}
		return nil
	}
}

func (j JApiCore) deleteByPathCore(operatingRoot string, delFunc func(num int64)) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		filename := req.Param("filename")
		if len(filename) == 0 {
			warning := "Failed to get FileName"
			return jutils.ProcessCustomHttpError("processUpload", warning, 400, w)
		}

		location := req.Param("location")
		if len(location) == 0 {
			warning := "Failed to get Location"
			return jutils.ProcessCustomHttpError("BasicDownloadByPathHandler", warning, 404, w)
		}

		cleanFilename := strings.ReplaceAll(filename, "/", "_")
		fmt.Println(cleanFilename)

		uniquePath := readUniquePath(req)
		if len(uniquePath) > 0 {
			operatingRoot += "/" + uniquePath
		}
		operatingRoot += "/" + location

		folder, err := j.FileIo.DownloadFolder(operatingRoot)
		if err != nil {
			jutils.ProcessHttpError("DeleteFile", err, 404, w)
			return err
		}

		deletionSize := folder.GetChildFiles()[cleanFilename].Size
		delFunc(deletionSize)

		err = j.FileIo.DeleteTargets([]string{cleanFilename}, folder)
		if err != nil {
			jutils.ProcessHttpError("DeleteFile", err, 500, w)
			return err
		}

		message := createJsonResponse("Deletion complete")
		condensedWriteJSON(w, message)
		return nil
	}
}

func (j JApiCore) ImportHandler() bunrouter.HandlerFunc {
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
			j.ScrapeQueue.Push(j.FileIo, w, &wg, operatingRoot, target, source)
		}

		wg.Wait()

		message := createJsonResponse("Import complete")
		condensedWriteJSON(w, message)
		return nil
	}
}

func (j JApiCore) UploadByPathHandler() bunrouter.HandlerFunc {
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

		fid := processUpload(w, j.FileIo, byteBuffer.Bytes(), head.Filename, operatingRoot, j.FileIoQueue)
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

func (j JApiCore) BasicDownloadFromBulkByPathHandler() bunrouter.HandlerFunc {
	operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_BULK_ROOT", "s/JAPI/Bulk")
	return j.downloadByPathCore(operatingRoot, func(num int64) {})
}

func (j JApiCore) BasicDownloadByPathHandler() bunrouter.HandlerFunc {
	operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_OP_ROOT", "s/JAPI")
	return j.downloadByPathCore(operatingRoot, func(num int64) {})
}

func (j JApiCore) BasicDeleteFromBulkByPathHandler() bunrouter.HandlerFunc {
	operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_BULK_ROOT", "s/JAPI/Bulk")
	return j.deleteByPathCore(operatingRoot, func(num int64) {})
}

func (j JApiCore) BasicDeleteByPathHandler() bunrouter.HandlerFunc {
	operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_OP_ROOT", "s/JAPI")
	return j.deleteByPathCore(operatingRoot, func(num int64) {})
}

func (j JApiCore) AdvancedDownloadFromBulkByPathHandler(reportFunc func(num int64)) bunrouter.HandlerFunc {
	operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_BULK_ROOT", "s/JAPI/Bulk")
	return j.downloadByPathCore(operatingRoot, reportFunc)
}

func (j JApiCore) AdvancedDownloadByPathHandler(reportFunc func(num int64)) bunrouter.HandlerFunc {
	operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_OP_ROOT", "s/JAPI")
	return j.downloadByPathCore(operatingRoot, reportFunc)
}

func (j JApiCore) AdvancedDeleteFromBulkByPathHandler(delFunc func(num int64)) bunrouter.HandlerFunc {
	operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_BULK_ROOT", "s/JAPI/Bulk")
	return j.deleteByPathCore(operatingRoot, delFunc)
}

func (j JApiCore) AdvancedDeleteByPathHandler(delFunc func(num int64)) bunrouter.HandlerFunc {
	operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_OP_ROOT", "s/JAPI")
	return j.deleteByPathCore(operatingRoot, delFunc)
}
