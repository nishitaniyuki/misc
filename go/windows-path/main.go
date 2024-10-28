package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	fmt.Println(filepath.Clean("a/b/c/d"))
	fmt.Println(filepath.Join("a/b/c", "d"))
}
