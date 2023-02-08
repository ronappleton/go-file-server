package filewatcher

import (
	"github.com/radovskyb/watcher"
	"github.com/ronappleton/golang-docker-base/filecache"
	"log"
	"os"
	"path/filepath"
	"time"
)

var fileWatcher *watcher.Watcher

func Watch(fileCache *filecache.FileCache) {
	fileWatcher = watcher.New()
	defer fileWatcher.Close()

	fileWatcher.SetMaxEvents(1)
	fileWatcher.FilterOps(watcher.Rename, watcher.Move, watcher.Create, watcher.Remove)

	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)

	if err := fileWatcher.AddRecursive(exPath + "/images"); err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			select {
			case event := <-fileWatcher.Event:
				fileCache.ProcessFileEvent(event)
			case err := <-fileWatcher.Error:
				log.Fatalln(err)
			case <-fileWatcher.Closed:
				return
			}
		}
	}()

	if err := fileWatcher.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
