package filewatcher

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/ronappleton/golang-docker-base/filecache"
	"github.com/ronappleton/golang-docker-base/files"
	"os"
	"path/filepath"
)

var watcher *fsnotify.Watcher
var fileCache *filecache.FileCache

func Watch(fCache *filecache.FileCache) {
	fileCache = fCache
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)

	if err := filepath.Walk(exPath+"/images", watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fileCache.ProcessFileEvent(event)
				break
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done
}

func filterFiles(data map[string]os.FileInfo) map[string]os.FileInfo {
	filtered := make(map[string]os.FileInfo, 0)

	for path, file := range data {
		if !file.IsDir() {
			filtered[path] = file
		}
	}

	return filtered
}

func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	} else {
		file := files.CreateFile(path)
		fileCache.Add <- &file
	}

	return nil
}
