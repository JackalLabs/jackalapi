package japicore

import (
	"net/http"
	"sync"

	"github.com/uptrace/bunrouter"

	"github.com/JackalLabs/jackalapi/jutils"
	"github.com/JackalLabs/jackalgo/handlers/file_io_handler"
	"github.com/JackalLabs/jackalgo/handlers/file_upload_handler"
)

func processUpload(w http.ResponseWriter, fileIo *file_io_handler.FileIoHandler, bytes []byte, cid string, pathSelect string, queue *FileIoQueue) string {
	fileUpload, err := file_upload_handler.TrackVirtualFile(bytes, cid, pathSelect)
	if err != nil {
		jutils.ProcessHttpError("TrackVirtualFile", err, 500, w)
		return ""
	}

	folder, msgs, err := fileIo.LoadNestedFolder(pathSelect)
	if err != nil {
		jutils.ProcessHttpError("LoadNestedFolder", err, 404, w)
		return ""
	}
	if len(msgs) > 0 {
		err = fileIo.SignAndBroadcast(msgs)
		if err != nil {
			jutils.ProcessHttpError("SignAndBroadcast", err, 500, w)
			return ""
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)

	m := queue.Push(fileUpload, folder, fileIo, &wg)

	wg.Wait()

	if m.Error() != nil {
		jutils.ProcessHttpError("UploadFailed", err, 500, w)
		return ""
	}

	return m.Fid()
}

func readUniquePath(req bunrouter.Request) string {
	uniquePath, ok := req.Context().Value(jutils.BasicKeyring.UseKey("ReqUniquePath")).(string)
	if !ok {
		return ""
	}
	if len(uniquePath) == 0 {
		return ""
	}
	return uniquePath
}
