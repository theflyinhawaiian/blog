package config

import (
	"os"
	"strings"
)

// ReadEnv returns the value for the given env var name.
// If a companion {name}_FILE env var is set, it reads the value from that file
// (the Docker secrets convention). Falls back to the plain env var otherwise.
func ReadEnv(name string) string {
	if path := os.Getenv(name + "_FILE"); path != "" {
		data, err := os.ReadFile(path)
		if err == nil {
			return strings.TrimRight(string(data), "\r\n")
		}
	}
	return os.Getenv(name)
}
