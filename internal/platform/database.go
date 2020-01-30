package platform

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-pg/pg"
	"io/ioutil"
	"time"

	_ "github.com/lib/pq"
)

// DBConfig is the required properties to use the database.
type DBConfig struct {
	Host       string
	ServerName string

	User     string
	Password string

	DisableTLS bool

	ServerCA      string
	ClientCert    string
	ClientKeyCert string
}

// OpenDB knows how to open a database connection based on the configuration.
func OpenDB(cfg DBConfig) (*pg.DB, error) {
	if cfg.DisableTLS {
		return pg.Connect(pgCfg(cfg, nil)), nil
	} else {
		// TLS enabled, load keys and configure TLS
		var serverCA []byte
		var err error
		if serverCA, err = ioutil.ReadFile(cfg.ServerCA); err != nil {
			return nil, err
		}

		cp := x509.NewCertPool()

		if ok := cp.AppendCertsFromPEM(serverCA); !ok {
			return nil, fmt.Errorf("unable to add cert for serverCA to cert pool. Config: %v", cfg)
		}
		clientCert := make([]tls.Certificate, 0, 1)
		certs, err := tls.LoadX509KeyPair(cfg.ClientCert, cfg.ClientKeyCert)
		if err != nil {
			return nil, fmt.Errorf("loadX509KeyPair error: %v", err)
		}
		clientCert = append(clientCert, certs)

		opt := pgCfg(cfg, &tls.Config{
			ServerName:         cfg.ServerName,
			Certificates:       clientCert,
			RootCAs:            cp,
			InsecureSkipVerify: false,
		})

		return pg.Connect(opt), nil
	}
}

func pgCfg(cfg DBConfig, c *tls.Config) *pg.Options {
	opt := &pg.Options{
		Addr:      fmt.Sprintf("%s:5432", cfg.Host),
		Database:  "postgres",
		User:      cfg.User,
		Password:  cfg.Password,
		TLSConfig: c,

		MaxRetries:      1,
		MinRetryBackoff: -1,

		DialTimeout:  30 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		PoolSize:           10,
		MaxConnAge:         10 * time.Second,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        10 * time.Second,
		IdleCheckFrequency: 100 * time.Millisecond,
	}
	return opt
}
