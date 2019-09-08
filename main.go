package main

import (
	"fmt"
	"goairmon/site"
	"goairmon/site/helper"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("failed to load .env")
	}

	serverCfg := &site.Config{
		Address:               helper.MustGetEnv("SERVER_ADDRESS"),
		AppCookieKey:          helper.MustGetEnv("APP_COOKIE_KEY"),
		CookieStoreEncryption: helper.MustGetEnv("COOKIE_STORE_ENCRYPTION"),
		StoragePath:           helper.MustGetEnv("STORAGE_PATH"),
	}

	server := site.NewSite(serverCfg)
	defer cleanup(server)

	server.Start()

	select {}
}

func cleanup(server *site.Site) {
	if err := server.Cleanup(); err != nil {
		fmt.Print(err)
	}
}
