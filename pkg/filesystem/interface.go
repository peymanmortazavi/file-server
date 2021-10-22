package filesystem

import (
	"io"
	"io/fs"
)

// Opener describes the ability to create a reader object.
type Opener interface {
	Open() (io.ReadCloser, error)
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
