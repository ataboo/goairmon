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

func BusinessRoot() string {
	return AppRoot() + "/business"
}

func SiteRoot() string {
	return AppRoot() + "/site"
}

func AppRoot() string {
	dir, _ := os.Getwd()
	for i := 0; i < 10; i++ {

		_, err := os.Stat(dir + "/site")
		if err == nil {
			return dir
		}

		dir = dir + "/.."
	}

	panic("failed to find site root")
}
