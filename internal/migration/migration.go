package migration

import (
	"context"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	timestampLayout = "2006-01-02T15:04Z"
	fileNameRe      = `\d+-[-A-Za-z]+\.sql$`
)

// Migrate executes all migration scripts found in the scriptsPath. In case any of scripts fails migration is stopped
// and error message will be returned.
func Migrate(ctx context.Context, db *pg.DB, scriptsPath string, logger *logrus.Entry) error {
	f, err := ioutil.ReadDir(scriptsPath)
	if err != nil {
		return err
	}

	timestampRe := regexp.MustCompile(`\d+`)
	fileNameRe := regexp.MustCompile(fileNameRe)
	m := make(map[time.Time]os.FileInfo)
	var keys []time.Time
	for _, fileInfo := range f {
		if strings.Contains(fileInfo.Name(), ".sql") {
			if err := validate(fileInfo.Name(), fileNameRe); err != nil {
				return err
			}

			s := timestampRe.FindAllString(fileInfo.Name(), 1)
			i, err := strconv.ParseInt(s[0], 10, 64)
			if err != nil {
				panic(err)
			}
			t := time.Unix(i, 0)
			m[t] = fileInfo
			keys = append(keys, t)
		}
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	// apply scripts in the right order
	for i := range keys {
		fileInfo := m[keys[i]]
		name := path.Join(scriptsPath, fileInfo.Name())
		logger.Infof("execute script: %v", name)
		b, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}

		if _, err := db.ExecContext(ctx, string(b)); err != nil {
			return err
		}
		logger.Infof("executed: %v", name)
	}

	return nil
}

func validate(fileName string, re *regexp.Regexp) error {
	if fileName == "" {
		return fmt.Errorf("file name cannot be empty")
	} else if !re.MatchString(fileName) {
		return fmt.Errorf("file name [%s] does not match the pattern. Example file name is '1580247785-add-column.sql'", fileName)
	}

	return nil
}
