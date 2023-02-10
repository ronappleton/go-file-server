package main

import (
	"flag"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

type config struct {
	port      string
	filesPath string
}

func (c *config) init(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	ex, err := os.Executable()

	if err != nil {
		panic("unable to get path of executable")
	}

	var (
		port      = flags.String("port", "80", "The port this server should run on")
		filesPath = flags.String("filesPath", filepath.Dir(ex)+"/files/", "Full path to file folder location")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	c.port = *port
	c.filesPath = *filesPath

	return nil
}

func main() {
	config := &config{}
	config.init(os.Args)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	app := fiber.New()
	app.Use(etag.New())
	app.Use(compress.New())

	app.Get("/*", func(c *fiber.Ctx) error {
		path := config.filesPath + "/" + c.Params("*")
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

	go func() {
		_ = <-done
		fmt.Println("Gracefully shutting down...")
		_ = app.ShutdownWithTimeout(time.Second * 5)
	}()
	
	if err := app.Listen(":" + config.port); err != nil {
		log.Panic(err)
	}
}
