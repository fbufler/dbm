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
	Databases    []database.Config `mapstructure:"databases"`
	DatabaseType string            `mapstructure:"database_type"`
}

func SetupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup dbm database tester",
		RunE:  runSetup,
	}
	cmd.Flags().StringSlice("databases", []string{}, "databases to test")
	cmd.Flags().String("database_type", "postgres", "database type to test")
	viper.BindPFlag("databases", cmd.Flags().Lookup("databases"))
	viper.BindPFlag("database_type", cmd.Flags().Lookup("database_type"))
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
		switch cfg.DatabaseType {
		case "sqlite":
			log.Debug().Msg("Using sqlite")
			dbs = append(dbs, database.NewSQLite(dbCfg))
		case "postgres":
			log.Debug().Msg("Using postgres")
			dbs = append(dbs, database.NewPostgres(dbCfg))
		}
	}
	tester := tester.New(tester.Config{
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
