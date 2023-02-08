package main

import (
	"github.com/ronappleton/golang-docker-base/filecache"
	"github.com/ronappleton/golang-docker-base/filewatcher"
)

func main() {
	fileCache := filecache.NewFileCache()
	go fileCache.Start()
	go filewatcher.Watch(fileCache)

	for {

	}
}
