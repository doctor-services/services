package env

import "os"

// GetEnvString gets value of environment variable with default value
func GetEnvString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
