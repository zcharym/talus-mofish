package main

import (
	"fmt"
	"os"

	"github.com/songwei.ma/talus-mofish/internal/database"
)

func main() {
	path, err := database.DefaultPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve database path: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(path)
}
