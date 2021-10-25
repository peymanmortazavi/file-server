package filesystem

type internalError struct {
	Message string
}

func (i internalError) Error() string {
	return i.Message
}

var (
	// FileAlreadyExists error for when a file already exists at a path.
	FileAlreadyExists = internalError{Message: "File already exists at the given path."}
)

// IsFileAlreadyExists returns if the error is the file already exists.
func IsFileAlreadyExists(err error) bool {
	if e, ok := err.(internalError); ok {
		return e == FileAlreadyExists
	}
	return false
}
