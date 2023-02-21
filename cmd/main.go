package main

import (
	"context"
	"expvar"
	"github.com/flow-lab/dlog"
	"github.com/flow-lab/flowdber/internal/migration"
	"github.com/flow-lab/flowdber/internal/platform"
	log "github.com/sirupsen/logrus"
	"os"
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
		AppName:      "flow-k8-sql",
		Level:        "debug",
		Version:      version,
		Commit:       short(commit),
		Build:        date,
		ReportCaller: true,
	})
	if err := run(logger); err != nil {
		logger.Fatalf("error: %v", err)
	}
}

func run(logger *log.Entry) error {
	serverName := os.Getenv("DB_SERVER_NAME")
	user := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	serverCA := os.Getenv("DB_SERVER_CA")
	clientCert := os.Getenv("DB_CLIENT_CERT")
	clientKeyCert := os.Getenv("DB_CLIENT_KEY_CERT")

	disableTLS := false
	if os.Getenv("DB_DISABLE_TLS") == "true" {
		disableTLS = true
	}

	c := platform.DBConfig{
		ServerName:    serverName,
		User:          user,
		Password:      pass,
		Host:          host,
		DisableTLS:    disableTLS,
		ServerCA:      serverCA,
		ClientCert:    clientCert,
		ClientKeyCert: clientKeyCert,
	}

	// db
	log.Println("database : initializing")
	db, err := platform.OpenDB(c)
	if err != nil {
		log.Printf("platform.OpenDB error: %v", err)
		return err
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Printf("database : closed with %s", err)
		} else {
			log.Printf("database : closed")
		}
	}()

	path := os.Getenv("DB_SQL_PATH")
	if path == "" {
		path = "/db"
	}
	err = migration.Migrate(context.Background(), db, path, logger)
	if err != nil {
		logger.Printf("error migrate: %v", err)
	}
	return err
}

func short(s string) string {
	if len(s) > 7 {
		return s[0:7]
	}
	return s
}
