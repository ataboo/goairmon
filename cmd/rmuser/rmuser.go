package main

import (
	"flag"
	"fmt"
	"goairmon/business/data/context"
	"goairmon/site/helper"

	"github.com/joho/godotenv"
)

func main() {
	envFilePath := flag.String("envpath", ".env", "path to .env file")
	userName := flag.String("username", "", "username to remove")

	flag.Parse()

	if err := godotenv.Load(*envFilePath); err != nil {
		panic("failed to load env file")
	}

	if *userName == "" {
		panic("username must be provided")
	}

	storagePath := helper.MustGetEnv("STORAGE_PATH")
	ctx := context.NewMemDbContext(&context.MemDbConfig{
		StoragePath: storagePath,
	})

	defer ctx.Close()

	user, err := ctx.FindUserByName(*userName)
	if err != nil {
		panic("user not found")
	}

	if err := ctx.DeleteUser(user.ID); err != nil {
		fmt.Println("failed to delete user", err)
		panic("exiting")
	}
}
