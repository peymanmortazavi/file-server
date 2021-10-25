package fshttp

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/peymanmortazavi/fs-server/pkg/filesystem"
)

// FileType demonstrates the type of a file item.
type FileType string

const (
	// RegularFile demonstrates regular file type.
	RegularFile FileType = "file"

	// DirType demonstrates dir file type.
	DirType FileType = "dir"
)

// FileItem represents an File.
type FileItem struct {
	Name       string      `json:"name,omitempty"`
	Type       FileType    `json:"type,omitempty"`
	Permission os.FileMode `json:"permission,omitempty"`
	Owner      string      `json:"owner,omitempty"`
	Size       int64       `json:"size,omitempty"`
	Data       string      `json:"data,omitempty"`
	Children   []FileItem  `json:"children,omitempty"`
}

// FileWriteRequest describes a file write request.
type FileWriteRequest struct {
	Data string `json:"data,omitempty"`
}

// CreateFileItemRequest represents a request to create a new file item.
type CreateFileItemRequest struct {
	FileWriteRequest `json:",inline"`
	Type             FileType `json:"type"`
}

// fileItemFromFSItem converts filesystem.Item to FileItem.
//
// This method only fails when populating data.
func fileItemFromFSItem(item filesystem.Item, populateData bool) (FileItem, error) {
	result := FileItem{
		Name:       item.Name,
		Size:       item.Size,
		Permission: item.Perm(),
		Owner:      item.Owner,
	}

	switch {
	case item.FileMode.IsDir():
		result.Type = DirType
		if item.Children != nil {
			result.Children = make([]FileItem, 0, len(item.Children))
			for _, child := range item.Children {
				childItem, err := fileItemFromFSItem(child, populateData)
				if err != nil {
					return result, err
				}
				result.Children = append(result.Children, childItem)
			}
		}
	case item.FileMode.IsRegular():
		result.Type = RegularFile
		if populateData && item.Opener != nil {
			builder := &strings.Builder{}
			file, err := item.Open(os.O_RDONLY)
			if err != nil {
				return result, fmt.Errorf("failed to read %s: %s", item.Name, err)
			}
			defer file.Close()
			io.Copy(builder, file)
			result.Data = builder.String()
		}
	}

	return result, nil
}
