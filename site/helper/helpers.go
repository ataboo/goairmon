package helper

import (
	"fmt"
	"os"
	"strconv"
)

func MustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("Failed to load .env value: %s", key))
	}

	return val
}

func MustGetEnvInt(key string) int {
	strVal := MustGetEnv(key)

	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert .env value: %s", key))
	}

	return intVal
}

func ResourceRoot() string {
	return AppRoot() + "/resources"
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

		_, err := os.Stat(dir + "/.env")
		if err == nil {
			return dir
		}

		dir = dir + "/.."
	}

	panic("failed to find site root")
}
