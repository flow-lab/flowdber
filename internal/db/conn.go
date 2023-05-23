package db

import (
	"database/sql"
	"fmt"
	utils "github.com/flow-lab/utils"
	_ "github.com/jackc/pgx/v4/stdlib"
	"os"
)

// ConnectTCPSocket initializes a TCP connection pool for a Cloud SQL
// instance of Postgres.
func ConnectTCPSocket() (*sql.DB, error) {
	var (
		dbUser          = utils.MustGetEnv("DB_USER") // e.g. 'my-db-user'
		dbPwd           = utils.MustGetEnv("DB_PASS") // e.g. 'my-db-password'
		dbTCPHost       = utils.MustGetEnv("DB_HOST") // e.g. '127.0.0.1' ('172.17.0.1' if deployed to GAE Flex)
		dbPort          = utils.MustGetEnv("DB_PORT") // e.g. '5432'
		dbName          = utils.MustGetEnv("DB_NAME") // e.g. 'my-database'
		sslRootCertName = utils.EnvOrDefault("DB_SSL_ROOT_CERT_NAME", "server-ca.pem")
		sslCertName     = utils.EnvOrDefault("DB_SSL_CERT_NAME", "server-ca.pem")
		sslKeyName      = utils.EnvOrDefault("DB_SSL_KEY_NAME", "client-key.pem")
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
		dbURI += fmt.Sprintf(" sslmode=require sslrootcert=%s/%s sslcert=%s/%s sslkey=%s/%s",
			dbCertPath, sslRootCertName, dbCertPath, sslCertName, dbCertPath, sslKeyName)
	}

	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	return dbPool, nil
}
