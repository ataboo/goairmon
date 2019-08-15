package main

import (
	"fmt"
	"goairmon/site"
)

func main() {
	server := site.NewSite()
	defer cleanup(server)

	server.Start()

	select {}
}

func cleanup(server *site.Site) {
	if err := server.Cleanup(); err != nil {
		fmt.Print(err)
	}
}
