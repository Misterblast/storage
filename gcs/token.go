package gcs

import (
	"context"
	"os"

	firebase "firebase.google.com/go"
	log "github.com/ghulammuzz/misterblast-storage/utils"
)

func CreateCustomToken(client *firebase.App) (string, error) {

	ctx := context.Background()
	c, err := client.Auth(ctx)
	if err != nil {
		log.Error("Failed to create new writer: %s", err.Error())
		return "", err
	}

	token, err := c.CustomToken(ctx, os.Getenv("STORAGE_SECRET"))
	if err != nil {
		log.Error("Failed to get custom token: %s", err.Error())
		return "", err
	}
	return token, nil
}
