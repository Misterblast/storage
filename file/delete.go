package file

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/ghulammuzz/misterblast-storage/utils"
)

var fileLock sync.Mutex

func Delete(key string) error {
	fileLock.Lock()
	defer fileLock.Unlock()

	destPath := filepath.Join(utils.UploadDir, key)

	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		return errors.New("file not found")
	}

	if err := os.Remove(destPath); err != nil {
		if os.IsNotExist(err) {
			return errors.New("file not found")
		}
		return errors.New("failed to delete file: " + err.Error())
	}

	return nil
}
