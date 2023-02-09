package filecache

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/ronappleton/golang-docker-base/files"
	"os"
)

type FileCache struct {
	Add       chan *files.File
	Remove    chan *files.File
	FindByKey chan *FileFinder
	FindByUrl chan *FileFinder
	Files     map[string]*files.File
}

type FileFinder struct {
	Find          string
	ReturnChannel chan *files.File
}

func NewFileFinder(find string) *FileFinder {
	return &FileFinder{
		Find:          find,
		ReturnChannel: make(chan *files.File),
	}
}

func (fileFinder *FileFinder) CloseReturnChannel() {
	close(fileFinder.ReturnChannel)
}

func NewFileCache() *FileCache {
	return &FileCache{
		Add:       make(chan *files.File),
		Remove:    make(chan *files.File),
		FindByKey: make(chan *FileFinder),
		FindByUrl: make(chan *FileFinder),
		Files:     make(map[string]*files.File),
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
		case fileFinder := <-fileCache.FindByKey:
			if file, ok := fileCache.Files[fileFinder.Find]; ok {
				fileFinder.ReturnChannel <- file
				break
			}

			fileFinder.ReturnChannel <- &files.File{}
			break
		case fileFinder := <-fileCache.FindByUrl:
			for _, file := range fileCache.Files {
				if file.UrlPath == fileFinder.Find {
					fileFinder.ReturnChannel <- file
					break
				}
				break
			}

			fileFinder.ReturnChannel <- &files.File{}
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
