package japicore

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/JackalLabs/jackalapi/jutils"
)

func httpGetFileRequest(w http.ResponseWriter, host string, path string) (bytes.Buffer, error) {
	var byteBuffer bytes.Buffer

	bytes.NewReader(byteBuffer.Bytes())

	url, err := url.Parse(host)
	if err != nil {
		jutils.ProcessHttpError("UrlParse", err, 500, w)
		return byteBuffer, err
	}

	url = url.JoinPath(path)
	// fmt.Println(url.String())

	innerReq, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		jutils.ProcessHttpError("CreateGetRequest", err, 500, w)
		return byteBuffer, err
	}

	res, err := http.DefaultClient.Do(innerReq)
	if err != nil {
		jutils.ProcessHttpError("UseGetRequest", err, 500, w)
		return byteBuffer, err
	}

	_, err = io.Copy(&byteBuffer, res.Body)
	if err != nil {
		jutils.ProcessHttpError("BufferCopy", err, 500, w)
		return byteBuffer, err
	}

	err = res.Body.Close()
	if err != nil {
		jutils.ProcessHttpError("BodyClose", err, 500, w)
		return byteBuffer, err
	}

	return byteBuffer, nil
}
