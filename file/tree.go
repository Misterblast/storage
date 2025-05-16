package file

import (
	"os"
	"path/filepath"

	"github.com/ghulammuzz/misterblast-storage/gcs"
	"github.com/gofiber/fiber/v2"
)

func buildLocalTree(path string) (*gcs.TreeNode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	node := &gcs.TreeNode{
		Name:   info.Name(),
		IsFile: !info.IsDir(),
	}

	if info.IsDir() {
		files, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			childPath := filepath.Join(path, f.Name())
			childNode, err := buildLocalTree(childPath)
			if err == nil {
				node.Children = append(node.Children, childNode)
			}
		}
	}

	return node, nil
}

func GetLocalTree(c *fiber.Ctx) error {
	rootPath := "./storage" // folder lokal proyek
	tree, err := buildLocalTree(rootPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(tree.Children)
}
