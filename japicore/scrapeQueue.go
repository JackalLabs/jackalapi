package japicore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/JackalLabs/jackalgo/handlers/file_io_handler"
	"github.com/JackalLabs/jutils"
)

type ScrapeQueue struct {
	scrapees []*scrapee
	fIQueue  *FileIoQueue
}

func NewScrapeQueue(fIQueue *FileIoQueue) *ScrapeQueue {
	q := ScrapeQueue{
		scrapees: make([]*scrapee, 0),
		fIQueue:  fIQueue,
	}
	return &q
}

type scrapee struct {
	fileIo   *file_io_handler.FileIoHandler
	w        http.ResponseWriter
	wg       *sync.WaitGroup
	err      error
	destPath string
	filename string
	host     string
}

func (m *scrapee) Error() error {
	return m.err
}

func (m *scrapee) Filename() string {
	return m.filename
}

func (m *scrapee) Host() string {
	return m.host
}

func (m *scrapee) LoadParts() (http.ResponseWriter, string) {
	return m.w, m.filename
}

func (q *ScrapeQueue) loadFIQueue() *FileIoQueue {
	return q.fIQueue
}

func (q *ScrapeQueue) Size() int {
	return len(q.scrapees)
}

func (q *ScrapeQueue) isEmpty() bool {
	return len(q.scrapees) == 0
}

func (q *ScrapeQueue) Push(fileIo *file_io_handler.FileIoHandler, w http.ResponseWriter, wg *sync.WaitGroup, destPath string, filename string, host string) *scrapee {
	m := scrapee{
		fileIo:   fileIo,
		w:        w,
		wg:       wg,
		destPath: destPath,
		filename: filename,
		host:     host,
	}

	q.scrapees = append(q.scrapees, &m)
	return &m
}

func (q *ScrapeQueue) pop() *scrapee {
	m := q.scrapees[0]
	q.scrapees = q.scrapees[1:]
	return m
}

func (q *ScrapeQueue) listenOnce() {
	if q.isEmpty() {
		return
	}

	scrapee := q.pop()
	w, filename := scrapee.LoadParts()

	byteBuffer, err := httpGetFileRequest(w, scrapee.Host(), filename)
	if err != nil {
		fmt.Println("ScrapeQueue.listenOnce() failure")
		return
	}

	fid := processUpload(w, scrapee.fileIo, byteBuffer.Bytes(), filename, scrapee.destPath, q.loadFIQueue())
	if len(fid) == 0 {
		warning := fmt.Sprintf("Failed to get FID for %s", filename)
		scrapee.err = jutils.ProcessCustomHttpError("processUpload", warning, 500, w)
		return
	}

	successfulUpload := UploadResponse{
		FID: fid,
	}
	err = json.NewEncoder(w).Encode(successfulUpload)
	if err != nil {
		jutils.ProcessHttpError("JSONSuccessEncode", err, 500, w)
		return
	}

	_, err = w.Write([]byte("uploadHandler"))
	if err != nil {
		jutils.ProcessError("UploadHandlerWrite", err)
	}

	scrapee.wg.Done()
}

func (q *ScrapeQueue) Listen() {
	for {
		q.listenOnce()
		time.Sleep(time.Second * 1)
	}
}
