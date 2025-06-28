package main

import (
	"flag"
	"fmt"
	mlog "log/slog"
	"os"
	"time"

	"github.com/ghulammuzz/misterblast-storage/database"
	"github.com/ghulammuzz/misterblast-storage/file"
	"github.com/ghulammuzz/misterblast-storage/gcs"
	"github.com/ghulammuzz/misterblast-storage/handler"
	"github.com/ghulammuzz/misterblast-storage/utils"
	log "github.com/ghulammuzz/misterblast-storage/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func init() {
	env := flag.String("env", "prod", "Environment for (stg/prod)")
	flag.Parse()

	if *env == "stg" {
		err := godotenv.Load("./.stg.env")
		if err != nil {
			mlog.Error("Error loading stg.env file ")
		}
		mlog.Info("Environment: staging (stg.env loaded)")
	} else {
		mlog.Info("Environment: production (using system environment variables)")
	}

	utils.Init()
	log.InitLogger("dev", false, "")
	// log.InitLogger("prod", true, "http://localhost:3100/loki/api/v1/push")

}

func main() {

	utils.InitStorage()
	db, err := database.InitPostgres()
	if err != nil {
		log.Error("Failed to initialize database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
		IdleTimeout:           30 * time.Second,
		BodyLimit:             15 * 1024 * 1024,
	})

	app.Use(utils.Cors())
	app.Use(logger.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	if err := godotenv.Load(); err != nil {
		log.Error("No .env file found")
	}

	if err := os.MkdirAll(utils.UploadDir, os.ModePerm); err != nil {
		log.Error("Failed to create upload directory: %v", err)
	}

	// app.Use(validateHeader)

	app.Static("/gcs", "./public")
	app.Static("/storage", "./storage")

	app.Post("/file", handler.Upload)
	app.Post("/file-admin", handler.UploadAdmin)
	app.Get("/file", handler.ServeImage)
	app.Get("/file/placeholder.png", func(c *fiber.Ctx) error {
		return c.SendFile("storage/placeholder.png")
	})

	// app.Post(d"/file", validateHeader, handler.Upload)
	// app.Get("/file", validateHeader, handler.ServeImage)
	// app.Get("/file/placeholder.png", validateHeader, func(c *fiber.Ctx) error {
	// 	return c.SendFile("storage/placeholder.png")
	// })

	app.Delete("/file", validateHeader, handler.Delete)
	app.Use("/file", func(c *fiber.Ctx) error { return validateHeader(c) })
	// app.Get("/token", handler.FirebaseToken)

	app.Get("/tree", gcs.GetStorageTree)
	app.Get("/local-tree", file.GetLocalTree)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Info("Server running on port:", port)

	go utils.StartPrometheusExporter()

	if err := app.Listen(fmt.Sprint(":", port)); err != nil {
		log.Error("Failed to start the server: %v", err)
	}
}

func validateHeader(c *fiber.Ctx) error {
	apiKey := os.Getenv("MISTERBLAST_KEY")
	clientKey := c.Get("MISTERBLAST_API_KEY")

	if clientKey == "" || clientKey != apiKey {
		log.Error("Unauthorized access attempt")
		return utils.JSON(c, 401, "Unauthorized: Invalid API key", nil)
	}
	return c.Next()
}
