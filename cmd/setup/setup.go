package setup

import (
	"github.com/fbufler/database-monitor/internal/tester"
	"github.com/fbufler/database-monitor/pkg/database"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type LocalCfg struct {
	Databases []database.Config `mapstructure:"databases"`
}

func SetupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup dbm database tester",
		RunE:  setup,
	}
	cmd.Flags().StringSlice("databases", []string{}, "databases to test")
	viper.BindPFlag("databases", cmd.Flags().Lookup("databases"))
	return cmd
}

func setup(cmd *cobra.Command, args []string) error {
	log.Info().Msg("Starting local")
	ctx := cmd.Context()
	LocalCfg := LocalCfg{}
	err := viper.Unmarshal(&LocalCfg)
	if err != nil {
		return err
	}
	log.Debug().Msgf("LocalCfg: %+v", LocalCfg)
	log.Debug().Msg("Initializing database tester")
	dbs := []database.Database{}
	for _, dbCfg := range LocalCfg.Databases {
		dbs = append(dbs, database.NewPostgres(dbCfg))
	}
	tester := tester.NewPostgres(tester.Config{
		Databases: dbs,
	})
	log.Info().Msg("Setup tester")
	err = tester.Setup(ctx)
	if err != nil {
		return err
	}
	log.Info().Msg("Tester setup complete")
	return nil
}
