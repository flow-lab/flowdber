package utils

import (
	"github.com/sirupsen/logrus"
	"os"
)

// GetEnvOrDefault returns the value of the environment variable k or defaultVal if it is not set.
func GetEnvOrDefault(k string, defaultVal string) string {
	v := os.Getenv(k)
	if v == "" {
		return defaultVal
	}
	return v
}

// MustGetEnv returns the value of the environment variable k or panics.
func MustGetEnv(k string, logger *logrus.Entry) string {
	v := os.Getenv(k)
	if v == "" && logger != nil {
		logger.Fatalf("Fatal Error in connect_tcp.go: %s environment variable not set.", k)
	}
	return v
}

// Short returns a shortened string.
func Short(s string) string {
	if len(s) > 7 {
		return s[0:7]
	}
	return s
}
