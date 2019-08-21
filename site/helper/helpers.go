package helper

import (
	"fmt"
	"os"
)

func MustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("Failed to load dotenv value: %s", key))
	}

	return val
}
