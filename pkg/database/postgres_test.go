package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresConnnectionString(t *testing.T) {
	cfg := Config{
		Host:              "localhost",
		Port:              5432,
		Username:          "testuser",
		Password:          "testpassword",
		Database:          "testdb",
		UseSSL:            true,
		SSLCertPath:       "/path/to/cert",
		SSLKeyPath:        "/path/to/key",
		SSLRootCertPath:   "/path/to/rootcert",
		ConnectionTimeout: 5,
	}
	expected := "host=localhost port=5432 user=testuser password=testpassword dbname=testdb sslmode=verify-full sslcert=/path/to/cert sslkey=/path/to/key sslrootcert=/path/to/rootcert connect_timeout=5"
	actual := cfg.postgresConnectionString()
	assert.Equal(t, expected, actual)
}

func TestPostgresConnnectionStringNoSSL(t *testing.T) {
	cfg := Config{
		Host:              "localhost",
		Port:              5432,
		Username:          "testuser",
		Password:          "testpassword",
		Database:          "testdb",
		UseSSL:            false,
		SSLCertPath:       "/path/to/cert",
		SSLKeyPath:        "/path/to/key",
		SSLRootCertPath:   "/path/to/rootcert",
		ConnectionTimeout: 5,
	}
	expected := "host=localhost port=5432 user=testuser password=testpassword dbname=testdb sslmode=disable connect_timeout=5"
	actual := cfg.postgresConnectionString()
	assert.Equal(t, expected, actual)
}

func TestPostgresConnnectionStringNoTimeout(t *testing.T) {
	cfg := Config{
		Host:              "localhost",
		Port:              5432,
		Username:          "testuser",
		Password:          "testpassword",
		Database:          "testdb",
		UseSSL:            true,
		SSLCertPath:       "/path/to/cert",
		SSLKeyPath:        "/path/to/key",
		SSLRootCertPath:   "/path/to/rootcert",
		ConnectionTimeout: 0,
	}
	expected := "host=localhost port=5432 user=testuser password=testpassword dbname=testdb sslmode=verify-full sslcert=/path/to/cert sslkey=/path/to/key sslrootcert=/path/to/rootcert connect_timeout=5"
	actual := cfg.postgresConnectionString()
	assert.Equal(t, expected, actual)
}
