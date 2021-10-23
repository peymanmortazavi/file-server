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

// DirManager is capable of managing a local directory.
// Vieweing and editing its content.
type DirManager struct {
	Root string
}

type fileOpener struct {
	path string
}

func (f fileOpener) Open() (io.ReadCloser, error) {
	file, err := os.Open(f.path)
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
