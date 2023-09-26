package japicore

import (
	"bytes"
	"encoding/json"
	"errors"
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

func Handler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		return nil
	}
}

func MethodNotAllowedHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		w.WriteHeader(http.StatusMethodNotAllowed)
		warning := fmt.Sprintf("%s method not availble for \"%s\"", req.URL.Path, req.Method)

		_, err := w.Write([]byte(warning))
		if err != nil {
			jutils.ProcessError("WWriteError for MethodNotAllowedHandler", err)
		}
		return nil
	}
}

func VersionHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		version := "v0.1.0"
		_, err := w.Write([]byte(version))
		if err != nil {
			jutils.ProcessError("WWriteError for VersionHandler", err)
		}
		return nil
	}
}

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

		_, err = w.Write([]byte("Import complete"))
		if err != nil {
			jutils.ProcessError("WWriteError for ImportHandler", err)
		}
		return nil
	}
}

func IpfsHandler(fileIo *file_io_handler.FileIoHandler, queue *FileIoQueue) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		var allBytes []byte

		operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_IPFS_ROOT", "s/JAPI/IPFS")
		gateway := jutils.LoadEnvVarOrFallback("JAPI_IPFS_GATEWAY", "https://ipfs.io/ipfs/")
		toClone := false
		cloneHeader := req.Header.Get("J-Clone-Ipfs")
		if strings.ToLower(cloneHeader) == "true" {
			toClone = true
		}

		id := req.Param("id")
		if len(id) == 0 {
			warning := "Failed to get IPFS CID"
			return jutils.ProcessCustomHttpError("processUpload", warning, 500, w)
		}

		cid := strings.ReplaceAll(id, "/", "_")

		handler, err := fileIo.DownloadFile(fmt.Sprintf("%s/%s", operatingRoot, cid))
		if err != nil {
			if !toClone {
				warning := "IPFS CID Not Found"
				w.WriteHeader(404)
				_, err := w.Write([]byte(warning))
				if err != nil {
					return err
				}
				return errors.New(strings.ToLower(warning))
			}

			byteBuffer, err := httpGetFileRequest(w, gateway, cid)
			if err != nil {
				jutils.ProcessHttpError("httpGetFileRequest", err, 404, w)
				return err
			}

			byteReader := bytes.NewReader(byteBuffer.Bytes())
			workingBytes := jutils.CloneBytes(byteReader)
			allBytes = jutils.CloneBytes(byteReader)

			fid := processUpload(w, fileIo, workingBytes, cid, operatingRoot, queue)
			if len(fid) == 0 {
				warning := "Failed to get FID post-upload"
				return jutils.ProcessCustomHttpError("IpfsHandler", warning, 500, w)
			}
		} else {
			allBytes = handler.GetFile().Buffer().Bytes()
		}
		_, err = w.Write(allBytes)
		if err != nil {
			jutils.ProcessError("WWriteError for IpfsHandler", err)
			return err
		}
		return nil
	}
}

func DownloadHandler(fileIo *file_io_handler.FileIoHandler) bunrouter.HandlerFunc {
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
			jutils.ProcessError("WWriteError for DownloadHandler", err)
		}
		return nil
	}
}

func UploadHandler(fileIo *file_io_handler.FileIoHandler, queue *FileIoQueue, subfolder string) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		var byteBuffer bytes.Buffer
		var wg sync.WaitGroup
		wg.Add(1)
		WorkingFileSize := 32 << 30

		operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_OP_ROOT", "s/JAPI")

		uniquePath := readUniquePath(req)
		if len(uniquePath) > 0 {
			operatingRoot += "/" + uniquePath
		}

		if len(subfolder) > 0 {
			operatingRoot += "/" + subfolder
		}

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

		_, err = w.Write([]byte("uploadHandler"))
		if err != nil {
			jutils.ProcessError("WWriteError for UploadHandler", err)
		}
		return nil
	}
}

func DeleteHandler(fileIo *file_io_handler.FileIoHandler, queue *FileIoQueue) bunrouter.HandlerFunc {
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

		return nil
	}
}
