package envvar

import (
	"fmt"
	"os"
)

func TryStringFromEnv(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("env variable %s not found", key)
	}

	return val, nil
}

func MustStringFromEnv(key string) string {
	val, err := TryStringFromEnv(key)
	if err != nil {
		panic(err)
	}
	return val
}
