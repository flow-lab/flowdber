package migration

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestMigrate(t *testing.T) {
	t.Run("Should validate file name", func(t *testing.T) {
		re := regexp.MustCompile(fileNameRe)
		assert.Nil(t, validate("1580247785-test-my-migration.sql", re))
		assert.Equal(t, fmt.Errorf("file name cannot be empty"), validate("", re))
	})
}
