package filesystem_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/peymanmortazavi/fs-server/pkg/filesystem"
)

const (
	aContent = "content of a.txt"
	bContent = "content of b.txt"
)

func createFileMap(result filesystem.Item) map[string]filesystem.Item {
	m := make(map[string]filesystem.Item, len(result.Children))
	for _, item := range result.Children {
		m[item.Name] = item
	}
	return m
}

func setupTestDir(t *testing.T) string {
	root, err := ioutil.TempDir("", "filesystem-dirmanager")
	if err != nil {
		t.Fatalf("failed to create test files.")
		return ""
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(root); err != nil {
			t.Logf("failed to clean up temporary directory: %s", root)
		}
	})
	return root
}

func createBasicDirStructure(root string) {
	ioutil.WriteFile(filepath.Join(root, "a.txt"), []byte(aContent), 0664)
	os.Mkdir(filepath.Join(root, "sub"), 0775)
	ioutil.WriteFile(filepath.Join(root, "sub", "b.txt"), []byte(bContent), 0664)
}

func TestView(t *testing.T) {
	root := setupTestDir(t)
	manager := filesystem.DirManager{root}

	item, err := manager.Get("/")
	if err != nil {
		t.Errorf("getting the root module failed: %s", err)
	}
	if !item.FileMode.IsDir() {
		t.Errorf("root module file mode is not dir")
	}
	if len(item.Children) != 0 {
		t.Errorf("expected no files in the root module but view returned: %d", len(item.Children))
	}

	type expectation struct {
		path     string
		dir      bool
		children []expectation
		data     string
	}

	createBasicDirStructure(root)

	cases := []expectation{
		{
			path: "sub",
			dir:  true,
			children: []expectation{
				{path: "b.txt", dir: false},
			},
		},
		{
			path: "sub/b.txt",
			dir:  false,
			data: bContent,
		},
		{
			path: "a.txt",
			dir:  false,
			data: aContent,
		},
	}

	for _, testCase := range cases {
		result, err := manager.Get(testCase.path)
		if err != nil {
			t.Errorf("manager.Get(%s) failed: %s", testCase.path, err)
		}
		if testCase.dir {
			if !result.FileMode.IsDir() {
				t.Errorf("expected dir for %s", testCase.path)
			}
			// create a map for the testCases to avoid sorting problems.
			childMap := createFileMap(result)
			for _, child := range testCase.children {
				info, ok := childMap[child.path]
				if !ok {
					t.Errorf("expected to find %s inside %s but it could not be found", child.path, testCase.path)
				}
				if child.dir && !info.FileMode.IsDir() {
					t.Errorf("expected %s inside %s to be listed as dir", child.path, testCase.path)
				}
				if !child.dir && !info.FileMode.IsRegular() {
					t.Errorf("expected %s inside %s to be listed as regular file", child.path, testCase.path)
				}
			}
		} else {
			// check file type.
			if !result.FileMode.IsRegular() {
				t.Errorf("expected regular file for %s", testCase.path)
			}
			// if data is provided, check data.
			if len(testCase.data) > 0 {
				reader, err := result.Open()
				if err != nil {
					t.Errorf("could not open file %s: %s", testCase.path, err)
				}
				buffer := &bytes.Buffer{}
				if _, err := io.Copy(buffer, reader); err != nil {
					t.Errorf("failed to get content of file %s: %s", testCase.path, err)
				}
				if testCase.data != buffer.String() {
					t.Errorf("incorrect data for %s: expected '%s' but got '%s'", testCase.path, testCase.data, buffer.String())
				}
			}
		}
	}
}
