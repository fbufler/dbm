package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type Postgres struct {
	Config     Config
	identifier string
	db         *sql.DB
}

func (c *Config) postgresConnectionString() string {
	log.Debug().Msg("Creating postgres connection string")
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", c.Host, c.Port, c.Username, c.Password, c.Database)
	if c.UseSSL {
		log.Debug().Msg("Using SSL")
		connectionString = fmt.Sprintf("%s sslmode=verify-full", connectionString)
		if c.SSLCertPath != "" {
			log.Debug().Msgf("Using SSL cert: %s", c.SSLCertPath)
			connectionString = fmt.Sprintf("%s sslcert=%s", connectionString, c.SSLCertPath)
		}
		if c.SSLKeyPath != "" {
			log.Debug().Msgf("Using SSL key: %s", c.SSLKeyPath)
			connectionString = fmt.Sprintf("%s sslkey=%s", connectionString, c.SSLKeyPath)
		}
		if c.SSLRootCertPath != "" {
			log.Debug().Msgf("Using SSL root cert: %s", c.SSLRootCertPath)
			connectionString = fmt.Sprintf("%s sslrootcert=%s", connectionString, c.SSLRootCertPath)
		}
	} else {
		log.Debug().Msg("Not using SSL")
		connectionString = fmt.Sprintf("%s sslmode=disable", connectionString)
	}
	if c.ConnectionTimeout == 0 {
		log.Debug().Msg("No connection timeout provided, using default connection timeout of 5 seconds")
		c.ConnectionTimeout = 5
	}
	connectionString = fmt.Sprintf("%s connect_timeout=%d", connectionString, c.ConnectionTimeout)
	return connectionString
}

func NewPostgres(cfg Config) Database {
	return &Postgres{
		Config:     cfg,
		identifier: fmt.Sprintf("%s:%d/%s", cfg.Host, cfg.Port, cfg.Database),
	}
}

func (p *Postgres) Connect() error {
	log.Debug().Msgf("%s: Connecting to postgres", p.identifier)
	db, err := sql.Open("postgres", p.Config.postgresConnectionString())
	if err != nil {
		return err
	}
	log.Debug().Msgf("%s: Connected to postgres", p.identifier)
	p.db = db
	return nil
}

func (p *Postgres) Close() error {
	log.Debug().Msgf("%s: Closing postgres connection", p.identifier)
	return p.db.Close()
}

func (p *Postgres) Identifier() string {
	return p.identifier
}

func (p *Postgres) Test(ctx context.Context) error {
	log.Debug().Msgf("%s: Testing postgres", p.identifier)
	if p.db == nil {
		log.Debug().Msgf("%s: No connection, connecting", p.identifier)
		err := p.Connect()
		if err != nil {
			return err
		}
		defer p.Close()
	}
	errChan := make(chan error)
	go func() {
		errChan <- p.TestWrite(ctx)
	}()
	go func() {
		errChan <- p.TestRead(ctx)
	}()
	returnCount := 0
	var testError error
	log.Debug().Msgf("%s: Waiting for test completion", p.identifier)
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%s: context terminated", p.identifier)
		case err := <-errChan:
			if err != nil {
				log.Debug().Msgf("%s: Received error: %s", p.identifier, err)
				if testError != nil {
					testError = fmt.Errorf("%s\n%s", testError, err)
				} else {
					testError = err
				}
			}
			returnCount++
		}
		if returnCount == 2 {
			log.Debug().Msgf("%s: All tests completed", p.identifier)
			return testError
		}
	}
}

func (p *Postgres) TestWrite(ctx context.Context) error {
	log.Debug().Msgf("%s, Testing postgres write", p.identifier)
	if p.db == nil {
		log.Debug().Msgf("%s: No connection, connecting", p.identifier)
		err := p.Connect()
		if err != nil {
			return err
		}
		defer p.Close()
	}
	log.Debug().Msgf("%s: Writing test data", p.identifier)
	_, err := p.db.ExecContext(ctx, "INSERT INTO test (test) VALUES ('test')")
	if err != nil {
		return err
	}
	log.Debug().Msgf("%s: Test data written", p.identifier)
	return nil
}

func (p *Postgres) TestRead(ctx context.Context) error {
	log.Debug().Msgf("%s: Testing postgres read", p.identifier)
	if p.db == nil {
		log.Debug().Msgf("%s: No connection, connecting", p.identifier)
		err := p.Connect()
		if err != nil {
			return err
		}
		defer p.Close()
	}
	log.Debug().Msgf("%s: Reading test data", p.identifier)
	rows, err := p.db.QueryContext(ctx, "SELECT * FROM test")
	if err != nil {
		return err
	}
	defer rows.Close()
	log.Debug().Msgf("%s: Test data read", p.identifier)
	return nil
}

func (p *Postgres) SetupTestTable(ctx context.Context) error {
	log.Debug().Msgf("%s: Setting up test table", p.identifier)
	if p.db == nil {
		log.Debug().Msgf("%s: No connection, connecting", p.identifier)
		err := p.Connect()
		if err != nil {
			return err
		}
		defer p.Close()
	}
	log.Debug().Msgf("%s: Dropping test table", p.identifier)
	_, err := p.db.ExecContext(ctx, "DROP TABLE IF EXISTS test")
	if err != nil {
		return err
	}
	log.Debug().Msgf("%s: Creating test table", p.identifier)
	_, err = p.db.ExecContext(ctx, "CREATE TABLE test (test varchar(255))")
	if err != nil {
		return err
	}
	log.Debug().Msgf("%s: Test table created", p.identifier)
	return nil
}
