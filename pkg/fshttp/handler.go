package fshttp

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/peymanmortazavi/fs-server/pkg/filesystem"
)

// Handler provides an HTTP interface to a file system handler.
type Handler struct {
	filesystem.Viewer
}

func writeError(writer http.ResponseWriter, e Error) {
	writer.WriteHeader(e.Status)
	if err := json.NewEncoder(writer).Encode(e); err != nil {
		log.Printf("failed to write error %s: %s", e, err)
	}
}

// Serve writes the response to the HTTP response.
func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var err error
	switch request.Method {
	case http.MethodGet:
		err = h.handleGet(writer, request)
	default:
		writeError(writer, methodNotAllowedError)
	}

	if err != nil {
		switch err := err.(type) {
		case Error:
			writeError(writer, err)
		default:
			writeError(writer, internalServerError)
		}
	}
}

func (h *Handler) handleGet(writer http.ResponseWriter, request *http.Request) error {
	// get the path
	path := strings.Trim(request.URL.Path, "/")
	item, err := h.Get(path)
	if err != nil {
		if os.IsNotExist(err) {
			return notFoundError
		}
		return err
	}
	query := request.URL.Query()
	populateData := item.FileMode.IsRegular() || query.Get("populateData") == "true"
	result, err := fileItemFromFSItem(item, populateData)
	if err != nil {
		log.Printf("failed to populate data for %s: %s", item.Name, err)
		return internalServerError
	}
	if err := json.NewEncoder(writer).Encode(result); err != nil {
		log.Printf("failed to write file item for %s: %s", path, err)
		return err
	}
	return nil
}
