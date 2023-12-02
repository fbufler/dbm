package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fbufler/database-monitor/cmd/local"
	"github.com/fbufler/database-monitor/cmd/serve"
	"github.com/fbufler/database-monitor/cmd/setup"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:     "dbm",
	Version: "0.0.1",
	Short:   "dbm is a database monitoring tool",
}

func initFlags() {
	rootCmd.PersistentFlags().Bool("debug", false, "sets log level to debug")
	rootCmd.PersistentFlags().String("logfile", "", "log file path")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))
}

func initViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/dbm/")
	viper.AddConfigPath("$HOME/.dbm")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("DBM")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No config file found, using defaults")
	}
}

func initLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	log.Logger = zerolog.New(output)
	debug := viper.GetBool("debug")
	file := viper.GetString("logfile")
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if file != "" {
		logFile, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		multi := zerolog.MultiLevelWriter(output, logFile)
		log.Logger = zerolog.New(multi).With().Timestamp().Logger()
	}
}

func main() {
	initViper()
	initFlags()
	initLogging()
	log.Debug().Msg("Starting dbm")
	log.Debug().Msg("Setting up root command")
	context, cancel := context.WithCancel(context.Background())
	rootCmd.SetContext(context)
	localCmd := local.LocalCommand()
	localCmd.SetContext(context)
	rootCmd.AddCommand(localCmd)
	setupCmd := setup.SetupCommand()
	setupCmd.SetContext(context)
	rootCmd.AddCommand(setupCmd)
	serveCmd := serve.ServeCommand()
	serveCmd.SetContext(context)
	rootCmd.AddCommand(serveCmd)
	log.Debug().Msg("Executing root command")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	go func() {
		err := rootCmd.Execute()
		if err != nil {
			log.Error().Err(err).Msg("Root command failed")
			sigchan <- syscall.SIGTERM
		}
	}()
	log.Debug().Msg("Waiting for termination signal")
	<-sigchan
	log.Debug().Msg("Received termination signal")
	cancel()
	log.Debug().Msg("Waiting for root command termination")
	time.Sleep(1 * time.Second)
	<-rootCmd.Context().Done()
	log.Debug().Msg("Root command terminated")
	log.Debug().Msg("Exiting dbm")
}
