package gcs

import (
	"context"

	firebase "firebase.google.com/go"
	"github.com/ghulammuzz/misterblast-storage/utils"
	log "github.com/ghulammuzz/misterblast-storage/utils"
)

func DeleteFileFromGCS(client *firebase.App, KeyFileName string) error {
	ctx := context.Background()

	c, err := client.Storage(ctx)
	if err != nil {
		log.Error("Failed to create new writer: %s", err.Error())
		return err
	}

	bucket, err := c.Bucket(utils.BucketName)
	if err != nil {
		log.Error("Failed to get bucket: %s", err.Error())
		return err
	}
	obj := bucket.Object(KeyFileName)
	if err := obj.Delete(ctx); err != nil {
		log.Error("Failed to delete file from GCS (%s): %s", KeyFileName, err.Error())
		return err
	}

	return nil
}
