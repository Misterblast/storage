package gcs

import (
	"context"
	"strings"

	"github.com/ghulammuzz/misterblast-storage/utils"
	"github.com/gofiber/fiber/v2"
)

type TreeNode struct {
	Name     string      `json:"name"`
	IsFile   bool        `json:"is_file"`
	Children []*TreeNode `json:"children,omitempty"`
}

func addToTree(root *TreeNode, path string) {
	parts := strings.Split(path, "/")
	current := root

	for i, part := range parts {
		var found *TreeNode
		for _, child := range current.Children {
			if child.Name == part {
				found = child
				break
			}
		}
		if found == nil {
			newNode := &TreeNode{
				Name:     part,
				IsFile:   i == len(parts)-1,
				Children: []*TreeNode{},
			}
			current.Children = append(current.Children, newNode)
			current = newNode
		} else {
			current = found
		}
	}
}

func GetStorageTree(c *fiber.Ctx) error {
	ctx := context.Background()

	app := utils.GCSClient
	if app == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "firebase storage not initialized",
		})
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	bucket, err := client.Bucket(utils.BucketName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	it := bucket.Objects(ctx, nil)
	root := &TreeNode{Name: "root", Children: []*TreeNode{}}

	for {
		obj, err := it.Next()
		if err != nil {
			break
		}
		addToTree(root, obj.Name)
	}

	return c.JSON(root.Children)
}
