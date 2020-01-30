package migration

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestMigrate(t *testing.T) {
	//run only with local db up and running
	//t.Run("Should read and apply all migration scripts to local db", func(t *testing.T) {
	//	cf := platform.DBConfig{
	//		Host:       "127.0.0.1",
	//		User:       "flowlab",
	//		DisableTLS: true,
	//	}
	//	db, err := platform.OpenDB(cf)
	//	assert.Nil(t, err)
	//
	//	err = Migrate(context.Background(), db, "testfiles", log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile|log.Ldate))
	//
	//	assert.Nil(t, err)
	//})

	t.Run("Should validate file name", func(t *testing.T) {
		re := regexp.MustCompile(fileNameRe)
		assert.Nil(t, validate("1580247785-test-my-migration.sql", re))
		assert.Equal(t, fmt.Errorf("file name cannot be empty"), validate("", re))
	})
}
