package jutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func CloneBytes(reader *bytes.Reader) []byte {
	var allBytes []byte
	_, err := reader.Read(allBytes)
	if err != nil {
		return nil
	}
	_, err = reader.Seek(0, 0)
	if err != nil {
		return nil
	}
	return allBytes
}

func CloneByteSlice(source []byte) ([]byte, []byte, error) {
	var firstSlice []byte
	var SecondSlice []byte
	byteReader := bytes.NewReader(source)

	_, err := byteReader.Read(firstSlice)
	if err != nil {
		return firstSlice, SecondSlice, err
	}

	_, err = byteReader.Seek(0, 0)
	if err != nil {
		return firstSlice, SecondSlice, err
	}

	_, err = byteReader.Read(SecondSlice)
	if err != nil {
		return firstSlice, SecondSlice, err
	}

	return firstSlice, SecondSlice, nil
}

func SimpleWriteJSON(w http.ResponseWriter, respVal interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(respVal); err != nil {
		ProcessHttpError("SimpleWriteJSON", err, 500, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	written, err := w.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("Written bytes: %d\n", written)
		ProcessError("SimpleWriteJSON", err)
	}
}
