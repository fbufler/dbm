package setup

import (
	"context"

	"github.com/fbufler/database-monitor/internal/tester"
	"github.com/fbufler/database-monitor/pkg/database"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SetupCfg struct {
	Databases []database.Config `mapstructure:"databases"`
}

func SetupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup dbm database tester",
		RunE:  runSetup,
	}
	cmd.Flags().StringSlice("databases", []string{}, "databases to test")
	viper.BindPFlag("databases", cmd.Flags().Lookup("databases"))
	return cmd
}

func runSetup(cmd *cobra.Command, args []string) error {
	log.Info().Msg("Starting local")
	ctx := cmd.Context()
	SetupCfg := SetupCfg{}
	err := viper.Unmarshal(&SetupCfg)
	if err != nil {
		return err
	}
	log.Debug().Msgf("SetupCfg: %+v", SetupCfg)
	return setup(&SetupCfg, ctx)
}

func setup(cfg *SetupCfg, ctx context.Context) error {
	log.Debug().Msg("Initializing database tester")
	dbs := []database.Database{}
	for _, dbCfg := range cfg.Databases {
		dbs = append(dbs, database.NewPostgres(dbCfg))
	}
	tester := tester.NewPostgres(tester.Config{
		Databases: dbs,
	})
	log.Info().Msg("Setup tester")
	err := tester.Setup(ctx)
	if err != nil {
		return err
	}
	log.Info().Msg("Tester setup complete")
	return nil
}
