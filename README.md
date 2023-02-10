# Golang File Server

Super simple file server

All this server does is serve static assets directly from the file system.

Features:

 - Server response caching
 - Etag for client caching
 - gzip compression
 - brotli compression
 - Automatic SSL (Lets Encrypt)

## Usage

First off `git clone https://github.com/ronappleton/go-file-server.git go-file-server`

Ensure go is installed (>1.16)
 
- In the root folder run `go build -buildvcs=false -o build/go-file-server`
- Copy the file from the build folder to your server somewhere.
- Create a folder next to the binary called `files` (I like to be original)
- Run `chmod +x go-file-server`
- Run `./go-file-server`

Go file server also accepts five optional arguments:

- `--port=80`
- `--storagePath=/your/folder/path`
- `--storageFolder=files`
- `--production=true`
- `--domain=www.example.com`
- `--email=admin@example.com`

These allow you to change the port the server serves from and allows you to set where your `files` folder lives,
as well as configure domain and email for ssl certificate issue, and production to enable ssl usage.

All there is left to do is point a domain at your server address and your in business!


## Why Iris?

I wanted a simple file server that just does its job as well as it can, and Iris is a web framework for go
that makes it real easy, adding the etag and compression functionality was easy.

Iris is the fastest HTTP2 compatible golang web framework.

[Iris](https://www.iris-go.com/)

[Github](https://github.com/kataras/iris)

![](https://github.com/kataras/server-benchmarks)
