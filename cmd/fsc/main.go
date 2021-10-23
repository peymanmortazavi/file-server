package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/peymanmortazavi/fs-server/pkg/fshttp"
)

func main() {
	host := flag.String("host", "0.0.0.0:6000", "the host to use for the API.")
	insecure := flag.Bool("insecure", false, "use insecure API.")
	raw := flag.Bool("raw", false, "get raw response.")

	flag.Parse()

	path := flag.Arg(0)

	scheme := "https"
	if *insecure {
		scheme = "http"
	}
	url, err := url.Parse(fmt.Sprintf("%s://%s/%s", scheme, *host, path))
	if err != nil {
		log.Fatalf("invalid host: %s", err)
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		log.Fatalf("failed to create request: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("request failed: %s", err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("request failed with status %d", resp.StatusCode)
		io.Copy(os.Stderr, resp.Body)
		return
	}

	if *raw {
		io.Copy(os.Stdout, resp.Body)
		return
	}

	var item fshttp.FileItem
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		log.Fatalf("failed to parse response as JSON: %s", err)
	}

	printItem(&item)
}

func printItem(item *fshttp.FileItem) {
	fmt.Printf("name: %s\n", item.Name)
	fmt.Printf("permission: %s\n", item.Permission)
	fmt.Printf("owner: %s\n", item.Owner)
	fmt.Printf("type: %s\n", item.Type)
	fmt.Printf("size (in bytes): %d\n", item.Size)
	switch item.Type {
	case fshttp.DirType:
		if len(item.Children) > 0 {
			fmt.Println()
			for _, child := range item.Children {
				fmt.Printf("%s  %-15s %-10d %5s   %s\n", child.Permission, child.Owner, child.Size, child.Type, child.Name)
			}
		}
	case fshttp.RegularFile:
		if len(item.Data) > 0 {
			fmt.Printf("\n%s\n", item.Data)
		}
	}
}
