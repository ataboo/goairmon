package main

import (
	"flag"
	"fmt"
	"goairmon/business/data/context"
	"goairmon/site/helper"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	envFilePath := flag.String("envpath", ".env", "path to .env file")
	userName := flag.String("username", "", "username to remove")

	flag.Parse()

	if err := godotenv.Load(*envFilePath); err != nil {
		fmt.Println("failed to load env file")
		os.Exit(1)
	}

	if *userName == "" {
		fmt.Println("username must be provided")
		os.Exit(1)
	}

	storagePath := helper.MustGetEnv("STORAGE_PATH")
	ctx := context.NewMemDbContext(&context.MemDbConfig{
		StoragePath: storagePath,
	})

	defer ctx.Close()

	user, err := ctx.FindUserByName(*userName)
	if err != nil {
		fmt.Println("user not found")
		os.Exit(1)
	}

	if err := ctx.DeleteUser(user.ID); err != nil {
		fmt.Println("failed to delete user", err)
		os.Exit(1)
	}

	if err := ctx.Save(); err != nil {
		fmt.Println("failed to save context", err)
		os.Exit(1)
	}

	fmt.Println("Success!")
}
