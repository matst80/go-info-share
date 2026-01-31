package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: cli <key> <value>")
		os.Exit(1)
	}

	key := os.Args[1]
	value := os.Args[2]

	url := fmt.Sprintf("http://localhost:8080/set?key=%s&value=%s", key, value)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", nil)
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