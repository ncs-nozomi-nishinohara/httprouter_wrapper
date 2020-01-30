package wrapper_utils

import "os"

func Getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return fallback
}

func Sqlenv(envname string, fallback string) string {
	var sql = "SQL_"
	var env = sql + envname
	return Getenv(env, fallback)
}
