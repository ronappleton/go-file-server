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
	port       string
	filesPath  string
	production bool
	domain     string
	email      string
}

func (c *config) init(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	ex, err := os.Executable()

	if err != nil {
		panic("unable to get path of executable")
	}

	var (
		port       = flags.String("port", "80", "The port this server should run on")
		filesPath  = flags.String("filesPath", filepath.Dir(ex), "Full path to file folder location")
		production = flags.Bool("production", false, "Run server in production mode")
		domain     = flags.String("domain", "", "The servers domain name")
		email      = flags.String("email", "", "The servers domain administrator")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	c.port = *port
	c.filesPath = *filesPath
	c.production = *production
	c.domain = *domain
	c.email = *email

	return nil
}

func main() {
	config := &config{}
	config.init(os.Args)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	app := iris.New()

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
		path := config.filesPath + "/files/" + ctx.Params().Get("root")
		file, err := os.Stat(path)
		fileHandle, err := os.Open(path)

		if err != nil {
			ctx.NotFound()
		}

		defer fileHandle.Close()

		ctx.CompressWriter(true)
		ctx.ServeContent(fileHandle, file.Name(), file.ModTime())
		ctx.ServeFile(path)
	})

	if config.production {
		app.Run(iris.AutoTLS(":"+config.port, config.domain, config.email))
	} else {
		if err := app.Listen(":" + config.port); err != nil {
			log.Panic(err)
		}
	}
	<-idleConnectionsClosed
}
