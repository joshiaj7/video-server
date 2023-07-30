package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

var (
	StoragePath = filepath.Join(".", "files")
)

type FileReader interface {
	GetName() string
	GetSize() int64
	GetFileMimeType() (string, error)
	Store() error
	Close() error
}

type fileReader struct {
	File   multipart.File
	Header *multipart.FileHeader

	fileMimeType string
	name         string
	size         int64
}

func NewFileReader(
	file multipart.File,
	header *multipart.FileHeader,
) *fileReader {
	return &fileReader{
		File:   file,
		Header: header,
		name:   "",
		size:   -1,
	}
}

func (f *fileReader) GetName() string {
	if f.name == "" {
		f.name = f.Header.Filename
	}
	return f.name
}

func (f *fileReader) GetSize() int64 {
	if f.size == -1 {
		f.size = f.Header.Size
	}
	return f.size
}

func (f *fileReader) GetFileMimeType() (string, error) {
	if f.fileMimeType == "" {
		mimeType, err := mimetype.DetectReader(f.File)
		if err != nil {
			panic(fmt.Sprintf("File Read Error: %s", err))
		}
		fileMimeType := mimeType.String()

		// back to 0 position
		_, err = f.File.Seek(0, io.SeekStart)
		if err != nil {
			panic(fmt.Sprintf("File Seek Error: %s", err))
		}

		f.fileMimeType = fileMimeType
	}

	return f.fileMimeType, nil
}

func (f *fileReader) Store() error {
	_ = os.MkdirAll(StoragePath, os.ModePerm)

	fullPath := StoragePath + "/" + f.GetName()
	osFile, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer osFile.Close()

	_, err = io.Copy(osFile, f.File)
	return err
}

func (f *fileReader) Close() error {
	closer, ok := interface{}(f.File).(io.Closer)
	if !ok {
		return nil
	}
	return closer.Close()
}
