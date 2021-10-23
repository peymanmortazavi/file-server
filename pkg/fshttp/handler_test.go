package fshttp_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/peymanmortazavi/fs-server/pkg/filesystem"
	"github.com/peymanmortazavi/fs-server/pkg/fshttp"
)

// an inefficient in-memory viewer implementation used for testing in this file.
type dummyViewer struct {
	Root filesystem.Item
}

// findItem returns a pointer to the found child file item if one exists.
func findItem(root *filesystem.Item, name string) *filesystem.Item {
	if root.Children == nil {
		return nil
	}
	for _, file := range root.Children {
		if file.Name == name {
			return &file
		}
	}
	return nil
}

func (d *dummyViewer) Get(path string) (filesystem.Item, error) {
	if path == "/" || path == "" {
		return d.Root, nil
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	cursor := &d.Root
	for _, part := range parts {
		file := findItem(cursor, part)
		if file == nil {
			return filesystem.Item{}, os.ErrNotExist
		}
		cursor = file
	}
	return *cursor, nil
}

type stringOpener string

type fakeFile struct {
	*strings.Reader
}

func (f fakeFile) Close() error {
	return nil
}

func (s stringOpener) Open() (io.ReadCloser, error) {
	return &fakeFile{strings.NewReader(string(s))}, nil
}

func mustMakeGETRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func TestView(t *testing.T) {
	root := filesystem.Item{
		Name:     "root",
		FileMode: os.ModeDir,
		Children: []filesystem.Item{
			{
				Name:     "sub",
				FileMode: os.ModeDir,
				Children: []filesystem.Item{
					{Name: "empty.txt", Opener: stringOpener("")},
				},
			},
			{
				Name:   "a.txt",
				Opener: stringOpener("a"),
			},
		},
	}

	viewer := &dummyViewer{root}
	handler := fshttp.Handler{viewer}

	testCases := []struct {
		request  *http.Request
		status   int
		children []string
		data     string
	}{
		{
			request:  mustMakeGETRequest("http://some.url.com/"),
			status:   200,
			children: []string{"sub", "a.txt"},
		},
		{
			request:  mustMakeGETRequest("http://some.url.com/sub"),
			status:   200,
			children: []string{"empty.txt"},
		},
		{
			request:  mustMakeGETRequest("http://some.url.com/sub?populateData=true"),
			status:   200,
			children: []string{"empty.txt"},
		},
		{
			request: mustMakeGETRequest("http://some.url.com/a.txt"),
			status:  200,
			data:    "a",
		},
		{
			request: mustMakeGETRequest("http://some.url.com/sub/empty.txt"),
			status:  200,
		},
		{
			request: mustMakeGETRequest("http://some.url.com/whooops"),
			status:  404,
		},
	}

	for _, testCase := range testCases {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, testCase.request)

		resp := recorder.Result()
		if testCase.status != resp.StatusCode {
			t.Errorf("unexpected status code for %s: expected %d, got %d",
				testCase.request.URL, testCase.status, resp.StatusCode)
		}
	}
}
