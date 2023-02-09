package main

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"log"
	"os"
	"path/filepath"
)

func main() {
	app := fiber.New()
	app.Use(etag.New())
	app.Use(compress.New())

	app.Get("/*", func(c *fiber.Ctx) error {
		ex, err := os.Executable()

		if err != nil {
			return fiber.ErrInternalServerError
		}

		exPath := filepath.Dir(ex)
		path := exPath + "/files/" + c.Params("*")
		mimeType, err := mimetype.DetectFile(path)

		if err != nil {
			return fiber.ErrInternalServerError
		}

		_, err = os.Stat(path)

		if err != nil {
			return fiber.ErrNotFound
		}

		data, err := os.ReadFile(path)

		if err != nil {
			return fiber.ErrInternalServerError
		}

		c.Set("Content-Type", mimeType.String())

		return c.Send(data)
	})

	log.Fatal(app.Listen(":80"))
}
