package japicore

import (
	"bytes"
	"fmt"
	"github.com/JackalLabs/jackalapi/jutils"
	"github.com/JackalLabs/jackalgo/handlers/file_io_handler"
	"github.com/uptrace/bunrouter"
	"net/http"
	"strings"
)

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
				return jutils.ProcessCustomHttpError("DownloadFile", warning, 404, w)
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

func DownloadByFidHandler(fileIo *file_io_handler.FileIoHandler) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		id := req.Param("id")
		if len(id) == 0 {
			warning := "Failed to get FileName"
			return jutils.ProcessCustomHttpError("DownloadByFidHandler", warning, 404, w)
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

func DeleteByFidHandler(fileIo *file_io_handler.FileIoHandler, queue *FileIoQueue) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		id := req.Param("id")
		if len(id) == 0 {
			warning := "Failed to get FileName"
			return jutils.ProcessCustomHttpError("processUpload", warning, 400, w)
		}

		fid := strings.ReplaceAll(id, "/", "_")
		fmt.Println(fid)

		// TODO - update after deletion by fid is added to jackalgo

		//folder, err := fileIo.DownloadFolder(queue.GetRoot("bulk"))
		//if err != nil {
		//	jutils.ProcessHttpError("DeleteFile", err, 404, w)
		//	return err
		//}
		//
		//err = fileIo.DeleteTargets([]string{fid}, folder)
		//if err != nil {
		//	jutils.ProcessHttpError("DeleteFile", err, 500, w)
		//	return err
		//}

		//message := createJsonResponse("Deletion complete")
		message := createJsonResponse("Deletion Not Implemented")
		condensedWriteJSON(w, message)
		return nil
	}
}
