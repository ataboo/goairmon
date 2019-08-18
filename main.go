package main

import (
	"fmt"
	"goairmon/site"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	server := site.NewSite()
	defer cleanup(server)

	err := godotenv.Load(".env")
	if err != nil {
		panic("failed to load .env")
	}

	serverCfg := &site.Config{
		Address:               mustGetEnv("SERVER_ADDRESS"),
		AppCookieKey:          mustGetEnv("APP_COOKIE_KEY"),
		CookieStoreEncryption: mustGetEnv("COOKIE_STORE_ENCRYPTION"),
	}

	server.Start(serverCfg)

	select {}
}

func cleanup(server *site.Site) {
	if err := server.Cleanup(); err != nil {
		fmt.Print(err)
	}
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("Failed to load dotenv value: %s", key))
	}

	return val
}
