package setup

import (
	"context"
	"os"
	"testing"

	"github.com/fbufler/database-monitor/pkg/database"
)

func TestSetup(t *testing.T) {
	cfg := SetupCfg{
		Databases: []database.Config{
			{
				FilePath: "test.db",
			},
		},
		DatabaseType: "sqlite",
	}
	ctx, cancel := context.WithCancel(context.Background())
	err := Setup(&cfg, ctx)
	if err != nil {
		t.Errorf("Setup failed: %s", err)
	}
	cancel()
	os.Remove("test.db")
}
