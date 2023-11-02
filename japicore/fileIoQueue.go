package japicore

import (
	"fmt"
	"sync"
	"time"

	"github.com/JackalLabs/jutils"

	"github.com/JackalLabs/jackalgo/handlers/file_io_handler"
	"github.com/JackalLabs/jackalgo/handlers/file_upload_handler"
	"github.com/JackalLabs/jackalgo/handlers/folder_handler"
)

type FileIoQueue struct {
	messages []*message
	roots    map[string]string
}

func NewFileIoQueue() *FileIoQueue {
	q := FileIoQueue{
		messages: make([]*message, 0),
		roots:    make(map[string]string),
	}

	ipfsRoot := jutils.LoadEnvVarOrFallback("JAPI_IPFS_ROOT", "s/JAPI/IPFS")
	q.PushRoot("ipfs", ipfsRoot)

	bulkRoot := jutils.LoadEnvVarOrFallback("JAPI_BULK_ROOT", "s/JAPI/Bulk")
	q.PushRoot("bulk", bulkRoot)

	return &q
}

type message struct {
	fileIo *file_io_handler.FileIoHandler
	upload *file_upload_handler.FileUploadHandler
	folder *folder_handler.FolderHandler
	wg     *sync.WaitGroup
	err    error
	fid    string
}

func (m *message) Error() error {
	return m.err
}

func (m *message) Fid() string {
	return m.fid
}

func (q *FileIoQueue) Size() int {
	return len(q.messages)
}

func (q *FileIoQueue) isEmpty() bool {
	return len(q.messages) == 0
}

func (q *FileIoQueue) Push(upload *file_upload_handler.FileUploadHandler, folder *folder_handler.FolderHandler, fileIo *file_io_handler.FileIoHandler, wg *sync.WaitGroup) *message {
	m := message{
		fileIo: fileIo,
		upload: upload,
		folder: folder,
		wg:     wg,
	}

	q.messages = append(q.messages, &m)
	return &m
}

func (q *FileIoQueue) PushRoot(root string, marker string) {
	q.roots[root] = marker
}

func (q *FileIoQueue) GetRoot(root string) string {
	return q.roots[root]
}

func (q *FileIoQueue) pop() *message {
	m := q.messages[0]
	q.messages = q.messages[1:]
	return m
}

func (q *FileIoQueue) listenOnce() {
	if q.isEmpty() {
		return
	}

	message := q.pop()

	_, fids, _, err := message.fileIo.StaggeredUploadFiles([]*file_upload_handler.FileUploadHandler{message.upload}, message.folder, false)

	fmt.Println(fids)

	message.err = err
	if len(fids) > 0 {
		message.fid = fids[0]
	}

	message.wg.Done()
}

func (q *FileIoQueue) Listen() {
	for {
		q.listenOnce()
		time.Sleep(time.Second * 5)
	}
}
