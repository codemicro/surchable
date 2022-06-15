package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("WebUI: %#v\n", os.Environ())
}
