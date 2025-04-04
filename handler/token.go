package handler

import (
	"github.com/ghulammuzz/misterblast-storage/gcs"
	"github.com/ghulammuzz/misterblast-storage/utils"
	"github.com/gofiber/fiber/v2"
)

func FirebaseToken(c *fiber.Ctx) error {

	token, err := gcs.CreateCustomToken(utils.GCSClient)
	if err != nil {
		return utils.JSON(c, 500, "Failed to create custom token", err.Error())
	}

	return utils.JSON(c, 200, "Custom token created successfully", token)
}
