package main

import (
	stdContext "context"
	"flag"
	"github.com/kataras/iris/v12"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

type config struct {
	port          string
	storagePath   string
	storageFolder string
	production    bool
	domain        string
	email         string
}

func (c *config) init(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	ex, err := os.Executable()

	if err != nil {
		panic("unable to get path of executable")
	}

	var (
		port          = flags.String("port", "80", "The port this server should run on")
		filesPath     = flags.String("storagePath", filepath.Dir(ex), "Full path to file folder location")
		storageFolder = flags.String("storageFolder", "files", "The outer most folder")
		production    = flags.Bool("production", false, "Run server in production mode")
		domain        = flags.String("domain", "", "The servers domain name")
		email         = flags.String("email", "", "The servers domain administrator")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	c.port = *port
	c.storagePath = *filesPath
	c.storageFolder = *storageFolder
	c.production = *production
	c.domain = *domain
	c.email = *email

	return nil
}

func main() {
	config := &config{}
	_ = config.init(os.Args)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	app := iris.New()
	app.Use(iris.Cache304(24 * time.Hour))
	app.Use(iris.Cache(720 * time.Hour))

	idleConnectionsClosed := make(chan struct{})
	iris.RegisterOnInterrupt(func() {
		timeout := 10 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		// close all hosts
		app.Shutdown(ctx)
		close(idleConnectionsClosed)
	})

	app.Get("{root:path}", func(ctx iris.Context) {
		path := config.storagePath + "/" + config.storageFolder + "/" + ctx.Params().Get("root")

		fileInfo, err := os.Stat(path)

		if err != nil {
			ctx.NotFound()
			return
		}

		file, _ := os.Open(path)
		defer file.Close()

		_ = ctx.CompressWriter(true)
		ctx.ServeContent(file, fileInfo.Name(), fileInfo.ModTime())
	})

	if config.production {
		if err := app.Run(iris.AutoTLS(":"+config.port, config.domain, config.email)); err != nil {
			log.Panic(err)
		}
	} else {
		if err := app.Listen(":" + config.port); err != nil {
			log.Panic(err)
		}
	}
	<-idleConnectionsClosed
}
