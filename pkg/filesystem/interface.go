package filesystem

import (
	"io"
	"io/fs"
)

// Opener describes the ability to create a reader object.
type Opener interface {
	Open(flag int) (io.ReadWriteCloser, error)
}

// Item describes a file system node.
type Item struct {
	fs.FileMode
	Name     string
	Owner    string
	Children []Item
	Size     int64
	Opener
}

// Viewer describes the ability to view or get items by path.
type Viewer interface {

	// Get retrieves one item by its path.
	Get(path string) (Item, error)
}

// Editor describes the ability to view or modify files.
type Editor interface {
	Viewer

	// CreateFile behaves like UNIX touch command, creating a file and returning it immediately.
	CreateFile(path string) (Item, error)

	// CreateDir creates a new directory, adding all sub paths in the way.
	CreateDir(path string) (Item, error)

	// Delete removes a file item at the given path.
	Delete(path string) error
}
