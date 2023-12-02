package serve

import (
	"github.com/fbufler/database-monitor/internal/service"
	"github.com/fbufler/database-monitor/internal/tester"
	"github.com/fbufler/database-monitor/pkg/database"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ServeCfg struct {
	Databases        []database.Config `mapstructure:"databases"`
	TestTimeout      int               `mapstructure:"test_timeout"`
	TestInterval     int               `mapstructure:"test_interval"`
	Port             int               `mapstructure:"port"`
	InvalidationTime int               `mapstructure:"invalidation_time"`
}

func ServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run dbm database tester service",
		RunE:  serve,
	}
	cmd.Flags().StringSlice("databases", []string{}, "databases to test")
	cmd.Flags().Int("test_timeout", 5, "test timeout in seconds")
	cmd.Flags().Int("test_interval", 5, "test interval in seconds")
	cmd.Flags().Int("port", 8080, "service port")
	cmd.Flags().Int("invalidation_time", 5, "invalidation time in seconds")
	viper.BindPFlag("databases", cmd.Flags().Lookup("databases"))
	viper.BindPFlag("test_timeout", cmd.Flags().Lookup("test_timeout"))
	viper.BindPFlag("test_interval", cmd.Flags().Lookup("test_interval"))
	viper.BindPFlag("port", cmd.Flags().Lookup("port"))
	viper.BindPFlag("invalidation_time", cmd.Flags().Lookup("invalidation_time"))
	return cmd
}

func serve(cmd *cobra.Command, args []string) error {
	log.Info().Msg("Starting local")
	ctx := cmd.Context()
	serveCfg := ServeCfg{}
	err := viper.Unmarshal(&serveCfg)
	if err != nil {
		return err
	}
	log.Debug().Msgf("LocalCfg: %+v", serveCfg)
	testTimeout := viper.GetInt("test_timeout")
	if testTimeout > 0 {
		serveCfg.TestTimeout = testTimeout
	}
	log.Debug().Msg("Initializing database tester")
	tester := tester.NewPostgres(tester.Config{
		Databases:    serveCfg.Databases,
		TestTimeout:  serveCfg.TestTimeout,
		TestInterval: serveCfg.TestInterval,
	})
	log.Info().Msg("Starting database tester")
	result := tester.Run(ctx)
	log.Info().Msg("Initializing service")
	router := mux.NewRouter()
	service := service.New(service.Config{
		Port:             serveCfg.Port,
		InvalidationTime: serveCfg.InvalidationTime,
	}, result, router)
	log.Info().Msg("Starting service")
	go service.Run(ctx)
	log.Debug().Msg("Waiting for context termination")
	<-ctx.Done()
	log.Info().Msg("Context terminated")
	return nil
}
