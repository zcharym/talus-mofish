package main

import (
	"fmt"
	"os"

	"github.com/songwei.ma/talus-mofish/backend/storage"
)

func main() {
	path, err := storage.DefaultPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve database path: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(path)
}
