package local

import (
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
		RunE:  local,
	}
	cmd.Flags().StringSlice("databases", []string{}, "databases to test")
	cmd.Flags().Int("test_timeout", 5, "test timeout in seconds")
	viper.BindPFlag("databases", cmd.Flags().Lookup("databases"))
	viper.BindPFlag("test_timeout", cmd.Flags().Lookup("test_timeout"))
	return cmd
}

func local(cmd *cobra.Command, args []string) error {
	log.Info().Msg("Starting local")
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
	log.Debug().Msg("Initializing database tester")
	dbs := []database.Database{}
	for _, dbCfg := range LocalCfg.Databases {
		dbs = append(dbs, database.NewPostgres(dbCfg))
	}
	tester := tester.NewPostgres(tester.Config{
		Databases:    dbs,
		TestTimeout:  LocalCfg.TestTimeout,
		TestInterval: LocalCfg.TestInterval,
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
