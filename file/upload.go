package file

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"

	"github.com/ghulammuzz/misterblast-storage/utils"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 4096)
	},
}

func Upload(fileHeader *multipart.FileHeader, key string) (string, error) {
	originalFileName := fileHeader.Filename
	destPath := filepath.Join(utils.UploadDir, key, originalFileName)
	dir := filepath.Dir(destPath)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", errors.New("failed to create directory")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", errors.New("failed to open uploaded file")
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return "", errors.New("failed to create destination file")
	}
	defer dst.Close()

	buffer := bufferPool.Get().([]byte)
	defer bufferPool.Put(buffer)

	for {
		n, err := src.Read(buffer)
		if n > 0 {
			if _, writeErr := dst.Write(buffer[:n]); writeErr != nil {
				return "", errors.New("failed to write file")
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", errors.New("failed to copy file")
		}
	}

	return destPath, nil
}
