package tester

import (
	"context"
	"time"

	"github.com/fbufler/database-monitor/pkg/database"
)

type Tester interface {
	Run(ctx context.Context) chan Result
	Setup(ctx context.Context) error
}

type Config struct {
	Databases    []database.Config `mapstructure:"databases"`
	TestTimeout  int               `mapstructure:"test_timeout"`
	TestInterval int               `mapstructure:"test_interval"`
}

type Result struct {
	Database       string        `json:"database"`
	Connectable    bool          `json:"connectable"`
	ConnectionTime time.Duration `json:"connection_time"`
	Writable       bool          `json:"writable"`
	WriteTime      time.Duration `json:"write_time"`
	Readable       bool          `json:"readable"`
	ReadTime       time.Duration `json:"read_time"`
	Timestamp      time.Time     `json:"timestamp"`
}
