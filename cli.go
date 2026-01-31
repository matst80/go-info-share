package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	var url string
	flag.StringVar(&url, "url", "", "Base URL of the info server")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("usage: cli [--url BASE_URL] <key> <value>")
		os.Exit(1)
	}

	key := args[0]
	value := args[1]

	if url == "" {
		url = os.Getenv("INFO_SERVER_URL")
		if url == "" {
			url = "http://localhost:8080"
		}
	}

	fullURL := fmt.Sprintf("%s/set?key=%s&value=%s", url, key, value)
	resp, err := http.Post(fullURL, "application/x-www-form-urlencoded", nil)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response:", err)
		os.Exit(1)
	}

	fmt.Println(string(body))
}