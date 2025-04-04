package utils

import (
	"context"
	"encoding/base64"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var GCSClient *firebase.App

func InitStorage() {
	ctx := context.Background()

	encodedCredentials := os.Getenv("FIREBASE_CREDENTIALS_BASE64")
	if encodedCredentials == "" {
		Error("error", "FIREBASE_CREDENTIALS_BASE64 not set")
		return
	}

	decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		Error("error decoding credentials", err.Error())
		return
	}

	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON(decodedCredentials))
	if err != nil {
		Error("error initial storage", err.Error())
		return
	}
	Info("success", "Firebase storage initialized")
	GCSClient = app
}
