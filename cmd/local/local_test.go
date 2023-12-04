package local

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/fbufler/database-monitor/cmd/setup"
	"github.com/fbufler/database-monitor/pkg/database"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

// Integration test for the local service.

func TestLocal(t *testing.T) {
	cfg := LocalCfg{
		Databases: []database.Config{
			{
				FilePath: "test.db",
			},
		},
		DatabaseType: "sqlite",
		TestTimeout:  1,
		TestInterval: 1,
	}
	ctx, cancel := context.WithCancel(context.Background())
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)
	setup.Setup(&setup.SetupCfg{
		Databases:    cfg.Databases,
		DatabaseType: cfg.DatabaseType,
	}, ctx)
	go local(&cfg, ctx)
	time.Sleep(1 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)

	bufS := buf.String()
	assert.Contains(t, bufS, "Starting local")
	assert.Contains(t, bufS, "Result: {Database:test.db")
	assert.Contains(t, bufS, "Context terminated")
	os.Remove("test.db")
}
