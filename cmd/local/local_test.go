package local

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/fbufler/database-monitor/pkg/database"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

// Integration test for the local service.

func TestLocal(t *testing.T) {
	cfg := LocalCfg{
		Databases: []database.Config{
			{
				Host:     "localhost",
				Port:     5432,
				Username: "postgres",
				Password: "postgres",
				Database: "postgres",
			},
		},
		TestTimeout:  1,
		TestInterval: 1,
	}
	ctx, cancel := context.WithCancel(context.Background())
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)
	go local(&cfg, ctx)
	time.Sleep(1 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)

	bufS := buf.String()
	assert.Contains(t, bufS, "Starting local")
	assert.Contains(t, bufS, "Result: {Database:localhost:5432/postgres")
	assert.Contains(t, bufS, "Context terminated")
}
