package japicore

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/JackalLabs/jackalapi/jutils"
	"github.com/uptrace/bunrouter"
)

func (j JApiCore) IpfsHandler() bunrouter.HandlerFunc {
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

		handler, err := j.FileIo.DownloadFile(fmt.Sprintf("%s/%s", operatingRoot, cid))
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

			clonedBytes1, clonedBytes2, err := jutils.CloneByteSlice(byteBuffer.Bytes())
			if err != nil {
				jutils.ProcessHttpError("httpGetFileRequest", err, 404, w)
				return err
			}
			allBytes = clonedBytes2

			fid := processUpload(w, j.FileIo, clonedBytes1, cid, operatingRoot, j.FileIoQueue)
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

func (j JApiCore) DownloadByFidHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		id := req.Param("id")
		if len(id) == 0 {
			warning := "Failed to get FileName"
			return jutils.ProcessCustomHttpError("DownloadByFidHandler", warning, 404, w)
		}
		fid := strings.ReplaceAll(id, "/", "_")

		handler, err := j.FileIo.DownloadFileFromFid(fid)
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

func (j JApiCore) DeleteByFidHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		id := req.Param("id")
		if len(id) == 0 {
			warning := "Failed to get FileName"
			return jutils.ProcessCustomHttpError("processUpload", warning, 400, w)
		}

		fid := strings.ReplaceAll(id, "/", "_")
		fmt.Println(fid)

		// TODO - update after deletion by fid is added to jackalgo

		//folder, err := j.FileIo.DownloadFolder(j.FileIoQueue.GetRoot("bulk"))
		//if err != nil {
		//	jutils.ProcessHttpError("DeleteFile", err, 404, w)
		//	return err
		//}
		//
		//err = j.FileIo.DeleteTargets([]string{fid}, folder)
		//if err != nil {
		//	jutils.ProcessHttpError("DeleteFile", err, 500, w)
		//	return err
		//}

		// message := createJsonResponse("Deletion complete")
		message := createJsonResponse("Deletion Not Implemented")
		condensedWriteJSON(w, message)
		return nil
	}
}
