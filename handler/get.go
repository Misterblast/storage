package handler

import (
	"fmt"
	"os"
	"path/filepath"

	db "github.com/ghulammuzz/misterblast-storage/database"
	"github.com/ghulammuzz/misterblast-storage/file"
	"github.com/ghulammuzz/misterblast-storage/utils"
	log "github.com/ghulammuzz/misterblast-storage/utils"

	"github.com/gofiber/fiber/v2"
)

func ServeImage(c *fiber.Ctx) error {
	key := c.Query("key")
	if key == "" {
		return utils.JSON(c, 400, "Key is required", nil)
	}

	filePath := filepath.Join(utils.UploadDir, key)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Error("[storage-server] File not found: %s", err.Error())
		log.Debug("Key Debug: %s", key)
		GCSurl, status, err := db.GetGCSURL(fmt.Sprintf("%s/%s?key=%s", os.Getenv("IMG_BASE_URL"), "file", key))
		log.Debug("GCS URL Debug: %s", GCSurl)

		if status == "pending" {
			log.Debug("File is still uploading, returning placeholder")
			return c.Redirect(fmt.Sprintf("%s/file/placeholder.png", os.Getenv("IMG_BASE_URL")))
		}

		if err != nil {
			log.Error("Failed to get GCS URL: %s", err.Error())
			return utils.JSON(c, 500, "Failed to get GCS URL", err.Error())
		}

		if err := file.Download(GCSurl, filePath); err != nil {
			log.Error("Failed to download file from GCS: %s", err.Error())
			return utils.JSON(c, 500, "Failed to download file from GCS", err.Error())
		}
		log.Debug("File downloaded from GCS: %s", filePath)

	}

	return c.SendFile(filePath)
}
