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
)

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
