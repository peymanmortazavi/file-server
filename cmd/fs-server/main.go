package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/peymanmortazavi/fs-server/pkg/filesystem"
	"github.com/peymanmortazavi/fs-server/pkg/fshttp"
)

func main() {
	rootDir := flag.String("root", "", "the root of the local path to serve.")
	addr := flag.String("addr", "0.0.0.0:6000", "the address to listen to (default: 0.0.0.0:6000)")

	flag.Parse()

	manager := filesystem.DirManager{Root: *rootDir}
	handler := &fshttp.Handler{Editor: manager}
	http.Handle("/", handler)

	log.Fatalln(http.ListenAndServe(*addr, nil))
}
