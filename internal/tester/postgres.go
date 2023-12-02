package tester

import (
	"context"
	"fmt"
	"time"

	"github.com/fbufler/database-monitor/pkg/database"
	"github.com/rs/zerolog/log"
)

type Postgres struct {
	results chan Result
	config  Config
}

func NewPostgres(config Config) Tester {
	return &Postgres{
		results: make(chan Result),
		config:  config,
	}
}

func (p *Postgres) run(ctx context.Context) {
	log.Info().Msg("Starting postgres tester")
	log.Debug().Msg("Initializing databases")
	dbs := make([]database.Database, len(p.config.Databases))
	for i, cfg := range p.config.Databases {
		dbs[i] = database.NewPostgres(cfg)
	}
	log.Debug().Msg("Starting database tests")
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("Received termination signal")
			log.Debug().Msg("Closing databases")
			for _, db := range dbs {
				db.Close()
			}
			log.Debug().Msg("Closing results channel")
			close(p.results)
			return
		default:
			for _, db := range dbs {
				go func(db database.Database) {
					result := Result{
						Database:    db.Identifier(),
						Connectable: false,
						Writable:    false,
						Readable:    false,
						Timestamp:   time.Now(),
					}
					connectionTime := time.Now()
					err := db.Connect()
					if err != nil {
						log.Error().Msgf("connecting to %s: %s", result.Database, err)
						p.results <- result
						return
					}
					result.ConnectionTime = time.Since(connectionTime)
					defer db.Close()
					result.Connectable = true
					readTime := time.Now()
					readCtx, readCancel := context.WithTimeout(ctx, time.Duration(p.config.TestTimeout)*time.Second)
					err = db.TestRead(readCtx)
					readCancel()
					if err != nil {
						log.Error().Msgf("reading from %s: %s", result.Database, err)
						p.results <- result
						return
					}
					result.ReadTime = time.Since(readTime)
					result.Readable = true
					writeTime := time.Now()
					writeCtx, writeCancel := context.WithTimeout(ctx, time.Duration(p.config.TestTimeout)*time.Second)
					err = db.TestWrite(writeCtx)
					writeCancel()
					if err != nil {
						log.Error().Msgf("writing to %s: %s", result.Database, err)
						p.results <- result
						return
					}
					result.WriteTime = time.Since(writeTime)
					result.Writable = true
					p.results <- result
				}(db)
			}
			time.Sleep(time.Duration(p.config.TestInterval) * time.Second)
		}
	}
}

func (p *Postgres) Run(ctx context.Context) chan Result {
	go p.run(ctx)
	return p.results
}

func (p *Postgres) Setup(ctx context.Context) error {
	var setupErrors []error
	for _, cfg := range p.config.Databases {
		db := database.NewPostgres(cfg)
		err := db.SetupTestTable(ctx)
		if err != nil {
			setupErrors = append(setupErrors, err)
		}
	}
	if len(setupErrors) > 0 {
		return fmt.Errorf("setting up databases: %v", setupErrors)
	}
	return nil
}
