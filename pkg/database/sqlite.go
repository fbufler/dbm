package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	Config     Config
	identifier string
	db         *sql.DB
}

func NewSQLite(cfg Config) Database {
	return &SQLite{
		Config:     cfg,
		identifier: cfg.FilePath,
	}
}

func (s *SQLite) Connect() error {
	db, err := sql.Open("sqlite3", s.Config.FilePath)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *SQLite) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *SQLite) Identifier() string {
	return s.identifier
}

func (s *SQLite) Test(ctx context.Context) error {
	if s.db == nil {
		err := s.Connect()
		if err != nil {
			return err
		}
		defer s.Close()
	}
	errChan := make(chan error)
	go func() {
		errChan <- s.TestWrite(ctx)
	}()
	go func() {
		errChan <- s.TestRead(ctx)
	}()
	returnCount := 0
	var testError error
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%s: context terminated", s.identifier)
		case err := <-errChan:
			if err != nil {
				if testError != nil {
					testError = fmt.Errorf("%s\n%s", testError, err)
				} else {
					testError = err
				}
			}
			returnCount++
			if returnCount == 2 {
				if testError != nil {
					return testError
				}
				return nil
			}
		}
	}
}

func (s *SQLite) TestWrite(ctx context.Context) error {
	if s.db == nil {
		err := s.Connect()
		if err != nil {
			return err
		}
		defer s.Close()
	}
	_, err := s.db.ExecContext(ctx, "INSERT INTO test (id, name) VALUES (1, 'test')")
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLite) TestRead(ctx context.Context) error {
	if s.db == nil {
		err := s.Connect()
		if err != nil {
			return err
		}
		defer s.Close()
	}
	_, err := s.db.ExecContext(ctx, "SELECT id, name FROM test WHERE id = 1")
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLite) SetupTestTable(ctx context.Context) error {
	if s.db == nil {
		err := s.Connect()
		if err != nil {
			return err
		}
		defer s.Close()
	}
	_, err := s.db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		return err
	}
	return nil
}
