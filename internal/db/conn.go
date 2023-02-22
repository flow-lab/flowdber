package db

import (
	"database/sql"
	"fmt"
	"github.com/flow-lab/flowdber/internal/utils"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
	"os"
)

// ConnectTCPSocket initializes a TCP connection pool for a Cloud SQL
// instance of Postgres.
func ConnectTCPSocket(logger *logrus.Entry) (*sql.DB, error) {
	var (
		dbUser    = utils.MustGetEnv("DB_USER", logger)  // e.g. 'my-db-user'
		dbPwd     = utils.GetEnvOrDefault("DB_PASS", "") // e.g. 'my-db-password'
		dbTCPHost = utils.MustGetEnv("DB_HOST", logger)  // e.g. '127.0.0.1' ('172.17.0.1' if deployed to GAE Flex)
		dbPort    = utils.MustGetEnv("DB_PORT", logger)  // e.g. '5432'
		dbName    = utils.MustGetEnv("DB_NAME", logger)  // e.g. 'my-database'
	)

	dbURI := fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s",
		dbTCPHost, dbUser, dbPwd, dbPort, dbName)

	// For deployments that connect directly to a Cloud SQL instance without
	// using the Cloud SQL Proxy, configuring SSL certificates will ensure the
	// connection is encrypted.
	if dbCertPath, ok := os.LookupEnv("DB_CERT_PATH"); ok { // e.g., '/path/to/db/certs'
		if dbCertPath[len(dbCertPath)-1:] == "/" {
			dbCertPath = dbCertPath[:len(dbCertPath)-1]
		}
		dbURI += fmt.Sprintf(" sslmode=require sslrootcert=%s/server-ca.pem sslcert=%s/client-cert.pem sslkey=%s/client-key.pem",
			dbCertPath, dbCertPath, dbCertPath)
	}

	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	return dbPool, nil
}
