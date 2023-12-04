package setup

import (
	"context"
	"testing"

	"github.com/fbufler/database-monitor/pkg/database"
)

func TestSetup(t *testing.T) {
	cfg := SetupCfg{
		Databases: []database.Config{
			{
				Host:     "localhost",
				Port:     5432,
				Username: "postgres",
				Password: "postgres",
				Database: "postgres",
			},
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	err := setup(&cfg, ctx)
	if err != nil {
		//TODO fix this if sqlite is available
		t.Skip(err)
	}
	cancel()
}
