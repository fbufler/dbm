package serve

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/fbufler/database-monitor/pkg/database"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestServe(t *testing.T) {
	cfg := ServeCfg{
		Port:             8081,
		InvalidationTime: 5,
		TestTimeout:      1,
		TestInterval:     1,
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
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)
	go serve(&cfg, ctx)
	time.Sleep(1 * time.Second)
	res, err := http.Get("http://localhost:8081/results")
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	cancel()
	_, err = http.Get("http://localhost:8081/results")
	assert.Error(t, err)
}
