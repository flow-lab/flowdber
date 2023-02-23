package main

import (
	"context"
	"expvar"
	"github.com/flow-lab/dlog"
	"github.com/flow-lab/flowdber/internal/db"
	"github.com/flow-lab/flowdber/internal/migration"
	utils "github.com/flow-lab/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	expvar.NewString("version").Set(version)
	expvar.NewString("commit").Set(commit)
	expvar.NewString("date").Set(date)

	logger := dlog.NewLogger(&dlog.Config{
		AppName:      "flowdber",
		Level:        utils.EnvOrDefault("LOG_LEVEL", "debug"),
		Version:      version,
		Commit:       utils.Short(commit),
		Build:        date,
		ReportCaller: true,
	})
	if err := run(logger); err != nil {
		logger.Fatalf("error: %v", err)
	}
}

func run(logger *log.Entry) error {
	// recover panics
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("panic: %v", r)
		}
	}()

	logger.Info("connect to db")
	dbConn, err := db.ConnectTCPSocket()
	if err != nil {
		return err
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			logger.Warnf("db close error: %s", err)
			return
		}
		logger.Infof("db connection closed")
	}()

	// ping connection
	if err := dbConn.Ping(); err != nil {
		return errors.Wrap(err, "db ping failed")
	}

	if err := migration.Migrate(
		context.Background(),
		dbConn,
		utils.EnvOrDefault("DB_SQL_PATH", "/db"),
		logger,
	); err != nil {
		return errors.Wrap(err, "migration failed")
	}

	return nil
}
