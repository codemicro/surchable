package main

import (
	"fmt"
	"os"

	db "github.com/codemicro/surchable/internal/libdb"
)

func main() {
	x, err := db.New()
	fmt.Println("db", x, err)
	conn, err := x.MakeConn()
	fmt.Println("dbc", conn, err)
	fmt.Printf("Coordinator: %#v\n", os.Environ())
}
