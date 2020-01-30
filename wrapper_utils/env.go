package wrapper_utils

import "os"

func Getenv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	return fallback
}

func Migration_env(envname string, fallback string) string {
	var MIGRATION = "MIGRATION_"
	var env = MIGRATION + envname
	return Getenv(env, fallback)
}
