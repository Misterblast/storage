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
	if client == nil {
		log.Error("GCS client is nil")
		return "", fmt.Errorf("GCS client is nil")
	}
	ctx := context.Background()

	log.Debug("[UploadFileToGCS] Uploading file to GCS: %s", file.Filename)

	src, err := file.Open()
	if err != nil {
		log.Error("Failed to open uploaded file: %s", err.Error())
		return "", err
	}
	defer src.Close()

	fileName := fmt.Sprintf("%s/%s", key, file.Filename)
	// wc := client.Bucket(utils.BucketName).Object(fileName).NewWriter(ctx)
	c, err := client.Storage(ctx)
	if err != nil {
		log.Error("Failed to create new writer: %s", err.Error())
		return "", err
	}

	bucket, err := c.Bucket(utils.BucketName)
	if err != nil {
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
		log.Error("Failed to write file to GCS (%s): %s", file.Filename, err.Error())
		return "", err
	}

	encodedFileName := url.QueryEscape(fileName)

	url := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media", utils.BucketName, encodedFileName)

	return url, nil
}
