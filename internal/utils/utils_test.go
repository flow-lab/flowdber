package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMustGetEnv(t *testing.T) {
	t.Run("TestMustGetEnv", func(t *testing.T) {
		// set env variable
		err := os.Setenv("TEST_KEY", "TEST_VALUE")
		assert.Nil(t, err)

		// test
		if MustGetEnv("TEST_KEY", nil) != "TEST_VALUE" {
			t.Errorf("MustGetEnv() = %s; want %s", MustGetEnv("TEST_KEY", nil), "TEST_VALUE")
		}

		// unset env variable
		err = os.Unsetenv("TEST_KEY")
		assert.Nil(t, err)
	})
}
