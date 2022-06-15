package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Coordinator: %#v\n", os.Environ())
}
