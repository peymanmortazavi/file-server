package filesystem

import (
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

const (
	defaultFlag = 0644
)

// DirManager is capable of managing a local directory.
// Vieweing and editing its content.
type DirManager struct {
	Root string
}

type fileOpener struct {
	path string
}

func (f fileOpener) Open(flag int) (io.ReadWriteCloser, error) {
	file, err := os.OpenFile(f.path, flag, defaultFlag)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func ownerName(info os.FileInfo) string {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid := stat.Uid
		userID := strconv.FormatUint(uint64(uid), 10)
		if user, err := user.LookupId(userID); err == nil {
			return user.Username
		}
	}
	return ""
}

// Get returns file item at the given path or returns an error.
//
// Error happens if file does not exist or the current user does not have permission
// to get (read).
func (d DirManager) Get(path string) (Item, error) {
	absolutePath := filepath.Join(d.Root, path)
	info, err := os.Stat(absolutePath)
	if err != nil {
		return Item{}, err
	}

	item := Item{FileMode: info.Mode(), Name: info.Name(), Size: info.Size(), Owner: ownerName(info)}

	if info.IsDir() {
		// create an item for the directory.
		files, err := ioutil.ReadDir(absolutePath)
		if err != nil {
			return item, err
		}
		item.Children = make([]Item, 0, len(files))
		for _, file := range files {
			child := Item{FileMode: file.Mode(), Name: file.Name(), Size: file.Size(), Owner: ownerName(file)}
			if file.Mode().IsRegular() {
				child.Opener = fileOpener{filepath.Join(absolutePath, file.Name())}
			}
			item.Children = append(item.Children, child)
		}
	} else {
		item.Opener = fileOpener{absolutePath}
	}
	return item, nil
}

// CreateFile creates a local file.
func (d DirManager) CreateFile(path string) (Item, error) {
	absolutePath := filepath.Join(d.Root, path)
	if _, err := os.Stat(absolutePath); err != nil {
		if os.IsNotExist(err) {
			if file, err := os.Create(absolutePath); err != nil {
				return Item{}, err
			} else {
				_ = file.Close()
				return d.Get(path)
			}
		}
		return Item{}, err
	}
	return Item{}, FileAlreadyExists
}

// CreateDir creates a local directory.
func (d DirManager) CreateDir(path string) (Item, error) {
	absolutePath := filepath.Join(d.Root, path)
	if err := os.MkdirAll(absolutePath, 0644); err != nil {
		return Item{}, err
	}
	return d.Get(path)
}

// Delete removes the file or directory completely.
func (d DirManager) Delete(path string) error {
	absolutePath := filepath.Join(d.Root, path)
	return os.RemoveAll(absolutePath)
}
