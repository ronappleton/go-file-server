package filecache

import (
	"errors"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gookit/goutil/dump"
	"github.com/radovskyb/watcher"
	"github.com/segmentio/fasthash/fnv1a"
	"os"
	"strconv"
	"strings"
)

var mimeTypes = []string{
	"image/jpeg",
	"image/bmp",
	"image/gif",
	"image/png",
	"image/vnd",
}

type File struct {
	key      string
	name     string
	data     string
	mimetype string
	filepath string
	urlPath  string
}

type FileCache struct {
	Add    chan *File
	Remove chan *File
	Files  map[string]*File
}

func NewFileCache() *FileCache {
	return &FileCache{
		Add:    make(chan *File),
		Remove: make(chan *File),
		Files:  make(map[string]*File),
	}
}

func (fileCache *FileCache) Start() {

	for {
		select {
		case file := <-fileCache.Add:
			dump.P(file)
			break
		case file := <-fileCache.Remove:
			dump.P(file)
			break
		}
	}
}

func (file *File) Add(fileCache *FileCache) {
	fileCache.Files[file.key] = file
}

func (file *File) Remove(fileCache *FileCache) {
	delete(fileCache.Files, file.key)
}

func (fileCache *FileCache) ProcessFileEvent(event watcher.Event) {
	mimeType, _ := mimetype.DetectFile(event.Path)
	if !contains(mimeTypes, mimeType.String()) {
		return
	}

	if event.Op == watcher.Create {
		file := createFile(event, mimeType.String())
		fileCache.Files[file.key] = &file
		dump.P(file.key, file.urlPath, file.mimetype, file.name, file.filepath)
	}

	if event.Op == watcher.Remove {

	}

	if event.Op == watcher.Rename || event.Op == watcher.Move {
		fmt.Println(event.Op)
		fmt.Println("event string:" + event.String())
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func createFile(event watcher.Event, mimeType string) File {
	data, _ := os.ReadFile(event.Path)

	return File{
		key:      strconv.FormatUint(fnv1a.HashString64(event.Path), 10),
		name:     event.Name(),
		data:     string(data),
		mimetype: mimeType,
		filepath: event.Path,
		urlPath:  getUrlPath(event.Path),
	}
}

func getUrlPath(path string) string {
	pathParts := strings.Split(path, "/")
	return trimmedPath(pathParts, "images")
}

func trimmedPath(pathParts []string, needle string) string {
	var idx = 0
	for index, str := range pathParts {
		if str == needle {
			idx = index
		}
	}

	return strings.Join(append(pathParts[:0], pathParts[idx+1:]...), "/")
}

func (fileCache *FileCache) GetFile(key string, urlPath string) (*File, error) {
	for _, file := range fileCache.Files {
		if file.key == key || file.urlPath == urlPath {
			return file, nil
		}
	}

	return nil, errors.New("unable to locate File")
}
