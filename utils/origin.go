package utils

import (
	"fmt"
	"os"
)

func CreateOriginURL(key string, types string, filename string) string {
	// http://localhost:3000/img?key=user/1/misterblast-2.png
	// http://localhost:3000/imguser/1
	return fmt.Sprintf("%s%s?key=%s/%s", os.Getenv("IMG_BASE_URL"), types, key, filename)
}
