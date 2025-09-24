package envvar

import (
	"fmt"
	"strconv"
)

func TryIntFromEnv(key string) (int, error) {
	s, err := TryStringFromEnv(key)
	if err != nil {
		return 0, err
	}

	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("cannot unparse %s env variable: %v", key, err)
	}

	return val, nil
}

func MustIntFromEnv(key string) int {
	val, err := TryIntFromEnv(key)
	if err != nil {
		panic(err)
	}
	return val
}
