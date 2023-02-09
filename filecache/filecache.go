package filecache

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/ronappleton/golang-docker-base/files"
	"os"
)

type FileCache struct {
	Add    chan *files.File
	Remove chan *files.File
	Find   chan *FileCacheReader
	Files  map[string]*files.File
}

type FileCacheReader struct {
	Key           string
	UrlPath       string
	ReturnChannel chan *files.File
}

func NewFileCacheReader(key string, urlPath string) *FileCacheReader {
	return &FileCacheReader{
		Key:           key,
		UrlPath:       urlPath,
		ReturnChannel: make(chan *files.File),
	}
}

func (fileCacheReader *FileCacheReader) CloseReturnChannel() {
	close(fileCacheReader.ReturnChannel)
}

func NewFileCache() *FileCache {
	return &FileCache{
		Add:    make(chan *files.File),
		Remove: make(chan *files.File),
		Find:   make(chan *FileCacheReader),
		Files:  make(map[string]*files.File),
	}
}

func (fileCache *FileCache) Start() {
	for {
		select {
		case file := <-fileCache.Add:
			fileCache.Files[file.Key] = file
			fmt.Println("File added: " + file.FilePath)
			break
		case file := <-fileCache.Remove:
			delete(fileCache.Files, file.Key)
			fmt.Println("File removed: " + file.FilePath)
			break
		case fileCacheReader := <-fileCache.Find:
			if fileCacheReader.Key != "" {
				if file, ok := fileCache.Files[fileCacheReader.Key]; ok {
					fileCacheReader.ReturnChannel <- file
					break
				}
			}

			if fileCacheReader.UrlPath != "" {
				for _, file := range fileCache.Files {
					if file.UrlPath == fileCacheReader.UrlPath {
						fileCacheReader.ReturnChannel <- file
						break
					}
				}
			}

			fileCacheReader.ReturnChannel <- &files.File{}
			break
		}
	}
}

func (fileCache *FileCache) ProcessFileEvent(event fsnotify.Event) {
	fmt.Println("FileWatcher event fired: " + event.String())
	_, err := files.GetMimeType(event.Name)
	if err != nil {
		fmt.Println("Unrecognised MimeType, returning...")
		return
	}

	// Create events are send when files are added and after the renaming of a file.
	if event.Has(fsnotify.Create) {
		fmt.Println("FileWatcher event: CREATE")
		file := files.CreateFile(event.Name)
		fileCache.Add <- &file
	}

	// When fsnotify sends a rename event, it actually sends the old filename and not the new one,
	// so we use that to remove from cache as it also sends a following create event for
	// the new filename.
	if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
		fmt.Println("FileWatcher event: REMOVE || RENAME")
		for _, file := range fileCache.Files {
			if file.FilePath == event.Name {
				fileCache.Remove <- file
			}
		}
	}

	// Chmod event can be sent on deleting a file (linux), so before we remove it from cache
	// we check the file has actually gone.
	if event.Has(fsnotify.Chmod) {
		fmt.Println("FileWatcher event: CHMOD")
		if _, err := os.Stat(event.Name); err != nil && errors.Is(err, os.ErrNotExist) {
			for _, file := range fileCache.Files {
				if file.FilePath == event.Name {
					fileCache.Remove <- file
				}
			}
		}
	}
}
