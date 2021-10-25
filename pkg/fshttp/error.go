package fshttp

import "net/http"

var (
	methodNotAllowedError = Error{
		Status:        http.StatusMethodNotAllowed,
		ID:            "method-not-allowed",
		UserMessage:   "Invalid request.",
		SystemMessage: "This method is not allowed for this endpoint.",
	}

	internalServerError = Error{
		Status:        http.StatusInternalServerError,
		ID:            "internal-server-error",
		UserMessage:   "oops! sorry something failed on our end.",
		SystemMessage: "unexpected server error occurred.",
	}

	notFoundError = Error{
		Status:        http.StatusNotFound,
		ID:            "not-found",
		UserMessage:   "oh no! no such file or directory.",
		SystemMessage: "no file element exists at the requested path or access is denied.",
	}

	fileExpected = Error{
		Status:        http.StatusBadRequest,
		ID:            "file-expected",
		UserMessage:   "this requested is only supported for files.",
		SystemMessage: "this request is only supported for files.",
	}

	jsonExpected = Error{
		Status:        http.StatusBadRequest,
		ID:            "bad-input",
		UserMessage:   "incorrect data format, only JSON is accepted.",
		SystemMessage: "could not parse request body as JSON.",
	}

	writeAccessDenied = Error{
		Status:        http.StatusForbidden,
		ID:            "write-access-denied",
		UserMessage:   "could not write to the requested file due to insufficient permission.",
		SystemMessage: "could not write to the requested file due to insufficient permission.",
	}

	deleteAccessDenied = Error{
		Status:        http.StatusForbidden,
		ID:            "delete-access-denied",
		UserMessage:   "you do not have permission to delete file or dir.",
		SystemMessage: "you do not have permission to delete file or dir.",
	}

	fileAlreadyExists = Error{
		Status:        http.StatusBadRequest,
		ID:            "file-already-exists",
		UserMessage:   "this file already exists at the given path, cannot create a new one.",
		SystemMessage: "this file already exists at the given path, cannot create a new one.",
	}
)

func newBadInputError(message string) Error {
	return Error{
		Status:        http.StatusBadRequest,
		ID:            "bad-input",
		UserMessage:   message,
		SystemMessage: message,
	}
}

// Error holds information about an error.
// TODO add http status here.
type Error struct {
	Status        int    `json:"-"`
	ID            string `json:"id,omitempty"`
	UserMessage   string `json:"user_message,omitempty"`
	SystemMessage string `json:"system_message,omitempty"`
}

func (e Error) Error() string {
	return e.UserMessage
}

func (e Error) String() string {
	return e.Error()
}
