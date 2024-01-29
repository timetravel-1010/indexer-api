package main

import (
	"fmt"
	"github.com/timetravel-1010/indexer-api/internal/server"
)

func main() {

	server := server.NewServer()

	fmt.Println("starting server...")
	err := server.ListenAndServe()
	if err != nil {
		panic("cannot start server")
	}
}
