package main

import (
	"flag"
	"fmt"
	"goairmon/business/data/context"
	"goairmon/business/data/models"
	"goairmon/site/helper"

	"github.com/joho/godotenv"
)

func main() {
	envFilePath := flag.String("envpath", ".env", "path to .env file")
	userName := flag.String("username", "", "username to add")
	password := flag.String("password", "", "password for user")

	flag.Parse()

	if err := godotenv.Load(*envFilePath); err != nil {
		fmt.Printf("failed to load env file %s", err)
	}

	if len(*userName) < 6 || len(*password) < 6 {
		fmt.Println("Username and Password must be atleast 6 characters")
		return
	}

	storagePath := helper.MustGetEnv("STORAGE_PATH")
	ctx := context.NewMemDbContext(&context.MemDbConfig{
		StoragePath: storagePath,
	})

	defer ctx.Close()

	user := models.User{
		Username: *userName,
	}
	if err := user.SetPassword(*password); err != nil {
		fmt.Printf("failed to set password: %s\n", err)
		return
	}

	if err := ctx.CreateOrUpdateUser(&user); err != nil {
		fmt.Printf("failed to create user: %s\n", err)
	}
}
