package jutils

import (
	"bytes"
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
