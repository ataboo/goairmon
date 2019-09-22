package main

import (
	"fmt"
	"goairmon/site"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("failed to load .env")
	}

	server := site.NewSite(site.EnvSiteConfig())
	defer cleanup(server)

	server.Start()

	select {}
}

func cleanup(server *site.Site) {
	if err := server.Cleanup(); err != nil {
		fmt.Print(err)
	}
}
