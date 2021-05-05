package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

const readHeaderBytes = 512

func processFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("os.Open(%s) failed: %s", file, err.Error())
	}
	defer f.Close()

	header := make([]byte, readHeaderBytes)
	if _, err := f.Read(header); err != nil {
		log.Fatalf("f.Read() failed: %s", err.Error())
	}

	s, err := f.Stat()
	if err != nil {
		log.Fatalf("os.Stat() failed: %s", err.Error())
	}

	ctype, _ := mimetype.Detect(header)

	fmt.Printf("File MIME type: %s\n", ctype)
	fmt.Printf("File size: %d\n", s.Size())
}

func processURL(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("http.Get(%s) failed: %s", url, err.Error())
	}
	defer resp.Body.Close()

	header := make([]byte, readHeaderBytes)
	if _, err := resp.Body.Read(header); err != nil {
		log.Fatalf("resp.Body.Read() failed: %s", err.Error())
	}

	ctype, _ := mimetype.Detect(header)

	fmt.Printf("Content-Length: %d\n", resp.ContentLength)
	fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))
	fmt.Printf("Detected Content-Type: %s\n", ctype)
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func isFile(s string) bool {
	fi, err := os.Stat(s)
	if err != nil {
		return false
	}

	return !fi.IsDir()
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("program takes exactly one argument that is either a URL or a file path")
	}
	fileOrURL := os.Args[1]

	if isURL(fileOrURL) {
		processURL(fileOrURL)
	} else if isFile(fileOrURL) {
		processFile(fileOrURL)
	} else {
		log.Fatalf("%s is neither a valid URL not a valid file", fileOrURL)
	}
}
