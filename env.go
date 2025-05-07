package zerodown

import "os"

type Env struct {
	key   string
	value string
}

var envs []Env

func SetEnv(key string, value string) {
	envs = append(envs, Env{key, value})
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
