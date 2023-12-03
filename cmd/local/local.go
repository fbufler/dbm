package local

import (
	"context"

	"github.com/fbufler/database-monitor/internal/tester"
	"github.com/fbufler/database-monitor/pkg/database"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type LocalCfg struct {
	Databases    []database.Config `mapstructure:"databases"`
	TestTimeout  int               `mapstructure:"test_timeout"`
	TestInterval int               `mapstructure:"test_interval"`
}

func LocalCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local",
		Short: "Run dbm database tester",
		RunE:  localRun,
	}
	cmd.Flags().StringSlice("databases", []string{}, "databases to test")
	cmd.Flags().Int("test_timeout", 5, "test timeout in seconds")
	viper.BindPFlag("databases", cmd.Flags().Lookup("databases"))
	viper.BindPFlag("test_timeout", cmd.Flags().Lookup("test_timeout"))
	return cmd
}

func localRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	LocalCfg := LocalCfg{}
	err := viper.Unmarshal(&LocalCfg)
	if err != nil {
		return err
	}
	log.Debug().Msgf("LocalCfg: %+v", LocalCfg)
	testTimeout := viper.GetInt("test_timeout")
	if testTimeout > 0 {
		LocalCfg.TestTimeout = testTimeout
	}
	return local(&LocalCfg, ctx)
}

func local(cfg *LocalCfg, ctx context.Context) error {
	log.Info().Msg("Starting local")

	log.Debug().Msg("Initializing database tester")
	dbs := []database.Database{}
	for _, dbCfg := range cfg.Databases {
		dbs = append(dbs, database.NewPostgres(dbCfg))
	}
	tester := tester.NewPostgres(tester.Config{
		Databases:    dbs,
		TestTimeout:  cfg.TestTimeout,
		TestInterval: cfg.TestInterval,
	})
	log.Info().Msg("Starting database tester")
	result := tester.Run(ctx)
	log.Debug().Msg("Handling results")
	for {
		select {
		case res := <-result:
			log.Info().Msgf("Result: %+v", res)
		case <-ctx.Done():
			log.Info().Msg("Context terminated")
			return nil
		}
	}
}
