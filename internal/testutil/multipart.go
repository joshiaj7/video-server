package testutil

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func RequestPayloadCreateFile(path string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("data", "sample.mp4")
	if err != nil {
		panic(err)
	}

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}

	err = writer.Close()
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "http://example.com/upload", body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}
