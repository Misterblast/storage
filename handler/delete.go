package handler

import (
	"fmt"
	"os"

	db "github.com/ghulammuzz/misterblast-storage/database"
	"github.com/ghulammuzz/misterblast-storage/file"
	"github.com/ghulammuzz/misterblast-storage/gcs"
	"github.com/ghulammuzz/misterblast-storage/utils"
	log "github.com/ghulammuzz/misterblast-storage/utils"
	"github.com/gofiber/fiber/v2"
)

func Delete(c *fiber.Ctx) error {
	key := c.Query("key")
	if key == "" {
		return utils.JSON(c, 400, "Key is required", nil)
	}

	tx, err := db.DBInstance.Begin()
	if err != nil {
		return utils.JSON(c, 500, "Failed to start transaction", err.Error())
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// keyFileName, err := file.GetFilename(key)
	// if err != nil {
	// 	return utils.JSON(c, 500, "Failed to get filename", err.Error())
	// }

	if err := file.Delete(key); err != nil {
		if err.Error() == "file not found" {
			return utils.JSON(c, fiber.StatusNotFound, "File not found", err.Error())
		}
		return utils.JSON(c, 500, "Failed to delete file", err.Error())
	}

	err = gcs.DeleteFileFromGCS(utils.GCSClient, key)
	if err != nil {
		return utils.JSON(c, 500, "Failed to delete file from GCS", err.Error())
	}

	originURL := fmt.Sprintf("%s/%s?key=%s", os.Getenv("IMG_BASE_URL"), "file", key)

	log.Debug("originURL: ", originURL)

	err = db.Delete(tx, originURL)
	if err != nil {
		return utils.JSON(c, 500, "Failed to delete file from database", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return utils.JSON(c, 500, "Failed to commit transaction", err.Error())
	}

	return utils.JSON(c, 200, "File deleted successfully", nil)
}
