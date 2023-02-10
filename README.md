# Golang File Server

Super simple file server

All this server does is serve static assets directly from the file system.

Features:

 - Etag for client caching
 - gzip compression
 - brotli compression

## Usage

First off `git clone https://github.com/ronappleton/go-file-server.git` go-file-server

Ensure go is installed (>1.16)
 
- In the root folder run `go build -buildvcs=false -o build/go-file-server`
- Copy the file from the build folder to your server somewhere.
- Create a folder next to the binary called `files` (I like to be original)
- Run `chmod +x go-file-server`
- Run `./go-file-server`

Go file server also accepts two optional arguments:

- `--port=80`
- `--filesPath=/your/folder/path`

These allow you to change the port the server serves from and allows you to set where your `files` folder lives

All there is left to do is point a domain at your server address and your in business!
## Why Fiber?

I wanted a simple file server that just does its job as well as it can, and Fiber is a web framework for go
that makes it real easy, adding the etag and compression middlewares was a snap, and it's built on top of
fasthttp, the images show why this is a good thing!

![](https://raw.githubusercontent.com/gofiber/docs/master/static/img/benchmark-pipeline.png)

![](https://raw.githubusercontent.com/gofiber/docs/master/static/img/benchmark_alloc.png)

