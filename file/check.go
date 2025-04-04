package file

import (
	"errors"
	"fmt"
	"os"

	"github.com/ghulammuzz/misterblast-storage/utils"
)

func CheckFile(fileName string) (bool, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}

func GetFilename(key string) (string, error) {
	files, err := os.ReadDir(utils.UploadDir + "/" + key)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", errors.New("file not found")
	}
	return fmt.Sprintf("%s%s", key, files[0].Name()), nil
}
