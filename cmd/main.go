package main

import (
	"context"
	"expvar"
	"fmt"
	"github.com/flow-lab/flow-k8-sql/internal/migration"
	"github.com/flow-lab/flow-k8-sql/internal/platform"
	"log"
	"os"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error : %s", err)
		os.Exit(1)
	}
}

func run() error {
	expvar.NewString("version").Set(version)
	expvar.NewString("commit").Set(commit)
	expvar.NewString("date").Set(date)
	log := log.New(os.Stdout, fmt.Sprintf("flow-k8-sql : (%s, %s) : ", version, short(commit)), log.LstdFlags|log.Lmicroseconds|log.Lshortfile|log.Ldate)

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
	err = migration.Migrate(context.Background(), db, path, log)
	if err != nil {
		log.Printf("error migrate: %v", err)
	}
	return err
}

func short(s string) string {
	if len(s) > 7 {
		return s[0:7]
	}
	return s
}
