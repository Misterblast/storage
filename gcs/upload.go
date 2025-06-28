package gcs

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"time"

	firebase "firebase.google.com/go"
	"github.com/ghulammuzz/misterblast-storage/utils"
	log "github.com/ghulammuzz/misterblast-storage/utils"
)

func UploadFileToGCS(client *firebase.App, file *multipart.FileHeader, key string) (string, error) {
	start := time.Now()
	destination := "gcs"

	defer func() {
		duration := time.Since(start).Seconds()
		utils.UploadDuration.WithLabelValues(destination).Observe(duration)
	}()

	if client == nil {
		utils.UploadErrorCounter.WithLabelValues(destination, "client_nil").Inc()
		log.Error("GCS client is nil")
		return "", fmt.Errorf("GCS client is nil")
	}

	ctx := context.Background()
	log.Debug("[UploadFileToGCS] Uploading file to GCS: %s", file.Filename)

	src, err := file.Open()
	if err != nil {
		utils.UploadErrorCounter.WithLabelValues(destination, "file_open_failed").Inc()
		log.Error("Failed to open uploaded file: %s", err.Error())
		return "", err
	}
	defer src.Close()

	fileSize := file.Size
	utils.UploadFileSize.WithLabelValues(destination).Observe(float64(fileSize))

	fileName := fmt.Sprintf("%s/%s", key, file.Filename)

	c, err := client.Storage(ctx)
	if err != nil {
		utils.UploadErrorCounter.WithLabelValues(destination, "storage_init_failed").Inc()
		log.Error("Failed to create storage client: %s", err.Error())
		return "", err
	}

	bucket, err := c.Bucket(utils.BucketName)
	if err != nil {
		utils.UploadErrorCounter.WithLabelValues(destination, "bucket_access_failed").Inc()
		log.Error("Failed to get bucket: %s", err.Error())
		return "", err
	}

	obj := bucket.Object(fileName)
	wc := obj.NewWriter(ctx)
	defer wc.Close()

	wc.ContentType = file.Header.Get("Content-Type")
	wc.Metadata = map[string]string{
		"uploaded_at": time.Now().Format(time.RFC3339),
		"key":         fileName,
	}

	if _, err := io.Copy(wc, src); err != nil {
		utils.UploadErrorCounter.WithLabelValues(destination, "gcs_write_failed").Inc()
		log.Error("Failed to write file to GCS (%s): %s", file.Filename, err.Error())
		return "", err
	}

	encodedFileName := url.QueryEscape(fileName)
	publicURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media", utils.BucketName, encodedFileName)

	utils.UploadRequestCounter.WithLabelValues(destination, "success").Inc()

	return publicURL, nil
}
