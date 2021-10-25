package fshttp

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/peymanmortazavi/fs-server/pkg/filesystem"
)

// Handler provides an HTTP interface to a file system handler.
type Handler struct {
	filesystem.Editor
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
	case http.MethodPost:
		err = h.handlePost(writer, request)
	case http.MethodPut:
		err = h.handlePut(writer, request)
	case http.MethodDelete:
		err = h.handleDelete(writer, request)
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

func writeToFile(opener filesystem.Opener, data string) error {
	file, err := opener.Open(os.O_CREATE | os.O_RDWR)
	if err != nil {
		if os.IsPermission(err) {
			return writeAccessDenied
		}
		return internalServerError
	}
	defer file.Close()
	if _, err := io.WriteString(file, data); err != nil {
		return internalServerError
	}
	return nil
}

func (h *Handler) handlePost(writer http.ResponseWriter, request *http.Request) error {
	path := strings.Trim(request.URL.Path, "/")
	if request.Body != nil {
		defer request.Body.Close()
	}
	var req CreateFileItemRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		return jsonExpected
	}
	switch req.Type {
	case RegularFile:
		item, err := h.CreateFile(path)
		if err != nil {
			if os.IsPermission(err) {
				return writeAccessDenied
			}
			if filesystem.IsFileAlreadyExists(err) {
				return fileAlreadyExists
			}
			log.Printf("failed to create file %s: %s", path, err)
			return internalServerError
		}
		if err := writeToFile(item, req.Data); err != nil {
			log.Printf("failed to write to file %s: %s", path, err)
			return err
		}
	case DirType:
		if _, err := h.CreateDir(path); err != nil {
			if os.IsPermission(err) {
				return writeAccessDenied
			}
			return internalServerError
		}
	default:
		return newBadInputError("invalid type, only file and dir are accepted.")
	}

	return nil
}

func (h *Handler) handlePut(writer http.ResponseWriter, request *http.Request) error {
	path := strings.Trim(request.URL.Path, "/")
	item, err := h.Get(path)
	if err != nil {
		if os.IsNotExist(err) {
			return notFoundError
		}
		return err
	}
	if item.IsDir() {
		return fileExpected
	}
	if request.Body != nil {
		defer request.Body.Close()
	}
	var req FileWriteRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		return jsonExpected
	}
	if err := writeToFile(item, req.Data); err != nil {
		log.Printf("failed to write to file %s: %s", path, err)
		return err
	}
	return nil
}

func (h *Handler) handleDelete(writer http.ResponseWriter, request *http.Request) error {
	path := strings.Trim(request.URL.Path, "/")
	if err := h.Delete(path); err != nil {
		switch {
		case os.IsNotExist(err):
			return notFoundError
		case os.IsPermission(err):
			return deleteAccessDenied
		}

		return err
	}
	return nil
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
