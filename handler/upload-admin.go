package handler

import (
	"fmt"
	"os"

	log "github.com/ghulammuzz/misterblast-storage/utils"

	db "github.com/ghulammuzz/misterblast-storage/database"
	"github.com/ghulammuzz/misterblast-storage/file"
	"github.com/ghulammuzz/misterblast-storage/gcs"
	"github.com/ghulammuzz/misterblast-storage/utils"

	"github.com/gofiber/fiber/v2"
)

func UploadAdmin(c *fiber.Ctx) error {
	files, err := c.FormFile("file")
	if err != nil {
		return utils.JSON(c, 400, "No file uploaded", nil)
	}

	const maxFileSize = 10 << 20
	if files.Size > maxFileSize {
		return utils.JSON(c, 400, "File size exceeds 10MB limit", nil)
	}

	key := c.FormValue("key")
	if key == "" {
		return utils.JSON(c, 400, "Key is required", nil)
	}

	path, err := file.Upload(files, key)
	if err != nil {
		return utils.JSON(c, 500, "Failed to upload file", nil)
	}

	originURL := utils.CreateOriginURL(key, "/file", files.Filename)
	publicURL := fmt.Sprintf("%s/%s?key=%s", os.Getenv("IMG_BASE_URL"), "file", path[8:])
	placeholderURL := fmt.Sprintf("%s/file/placeholder.png", os.Getenv("IMG_BASE_URL"))

	log.Debug("placeholderURL: ", placeholderURL)

	tx, err := db.DBInstance.Begin()
	if err != nil {
		log.Error("Failed to start DB transaction: %v", err)
		return utils.JSON(c, 500, "Database transaction errorpath[8:]", nil)
	}

	err = db.CreateOrUpdate(tx, key, originURL, placeholderURL, "pending")
	if err != nil {
		log.Error("DB Update failed: %v", err)
		tx.Rollback()
		return utils.JSON(c, 500, "Failed to update database", nil)
	}
	tx.Commit()

	go func() {
		gcsURL, err := gcs.UploadFileToGCS(utils.GCSClient, files, key)
		if err != nil {
			log.Error("GCS Upload failed: %v", err)
			return
		}
		log.Debug("GCS URL: %s", gcsURL)

		tx, err := db.DBInstance.Begin()
		if err != nil {
			log.Error("Failed to start DB transaction: %v", err)
			return
		}

		err = db.CreateOrUpdate(tx, key, originURL, gcsURL, "completed")
		if err != nil {
			log.Error("DB Update failed: %v", err)
			tx.Rollback()
			return
		}
		log.Debug("DB Update completed")

		err = tx.Commit()
		if err != nil {
			log.Error("DB Commit failed: %v", err)
		}
	}()

	return utils.JSON(c, 200, "File uploaded successfully", fiber.Map{
		"url": publicURL,
	})
}
