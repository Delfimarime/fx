package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/delfimarime/fx/config"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
)

const (
	App = "switch"
)

func RunHttpServerSpf13CobraCommand(seq ...func(*Opts)) {
	rootCmd := NewHttpServerSpf13CobraCommand(seq...)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func RunSpf13CobraCommand(f func(Opts) *cobra.Command, seq ...func(*Opts)) {
	opts := NewOpts(seq...)
	rootCmd := f(opts)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewHttpServerSpf13CobraCommand(seq ...func(*Opts)) *cobra.Command {
	opts := NewOpts(seq...)
	rootCmd := NewRootSpf13CobraCommand(opts)
	startupCmd := NewStartServerSpf13CobraCommand(opts)
	rootCmd.AddCommand(startupCmd)
	return rootCmd
}

func NewRootSpf13CobraCommand(opts Opts) *cobra.Command {
	return &cobra.Command{
		Use:   App,
		Short: fmt.Sprintf("HTTP Server that exposes %s REST API", opts.api),
		RunE: func(cmd *cobra.Command, args []string) error {
			binary, err := os.ReadFile("package.json")
			if err != nil {
				return err
			}
			packageInformation := make(map[string]interface{})
			err = json.Unmarshal(binary, &packageInformation)
			if err != nil {
				return err
			}
			version := packageInformation["version"]
			if version == "" {
				version = "N/A"
			}
			cmd.Println(fmt.Sprintf("%s\nversion:%s\nREST API:%s", App, version, opts.api))
			return nil
		},
	}
}

func NewStartServerSpf13CobraCommand(opts Opts) *cobra.Command {
	startupCmd := &cobra.Command{
		Use:   "up",
		Short: fmt.Sprintf("Starts %s as an HTTP Server to expose %s REST API.", App, opts.api),
		Args:  cobra.MaximumNArgs(1),
		Run:   ExecuteStartServerSpf13CobraCommand(opts),
	}
	// Add flags to the startup command
	startupCmd.PersistentFlags().String("log-level", "INFO", "Set the logging level (e.g., DEBUG, INFO)")
	startupCmd.PersistentFlags().Int("server-port", 8080, "Port on which the server runs")
	startupCmd.PersistentFlags().String("server-mode", "release", "The mode the server is running")
	return startupCmd
}

func ExecuteStartServerSpf13CobraCommand(opts Opts) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		configuration, err := getSpf13CobraTerminalConfiguration(cmd, args)
		if err != nil {
			panic(err)
		}
		zapLevel, err := config.ToZapLogLevel(configuration.Logging.Level)
		if err != nil {
			panic(err)
		}
		// CONFIGURE LOGGING
		logConfiguration := zap.NewProductionConfig()
		logConfiguration.Level.SetLevel(zapLevel)
		logger, err := logConfiguration.Build()
		if err != nil {
			panic(err)
		}
		defer logger.Sync()
		zap.ReplaceGlobals(logger)
		startOpts := make([]fx.Option, 0)
		for _, each := range opts.factories {
			startOpts = append(startOpts, each(configuration.Config))
		}
		startOpts = append(startOpts, opts.options...)
		startOpts = append(startOpts, fx.Provide(func() config.Terminal {
			return configuration
		}, func(c config.Terminal) config.Config {
			return c.Config
		}))
		fx.New(startOpts...).Run()
	}
}

func getSpf13CobraTerminalConfiguration(cmd *cobra.Command, args []string) (config.Terminal, error) {
	configURI := ""
	if len(args) > 0 {
		configURI = args[0]
	}
	if configURI == "" {
		value := os.Getenv(strings.ToUpper(fmt.Sprintf("%s_HOME", App)))
		if value == "" {
			value = "."
		}
		configURI = fmt.Sprintf("%s/config.yaml", value)
	}
	configLoader := config.NewConfigLoader(configURI)
	target := config.Terminal{}
	err := configLoader.ReadConfiguration(&target)
	if err != nil {
		return config.Terminal{}, err
	}
	valueOf, _ := cmd.Flags().GetString("log-level")
	if valueOf != "" {
		target.Logging.Level = valueOf
	}
	valueOf, _ = cmd.Flags().GetString("server-port")
	if valueOf != "" {
		serverPort, err := strconv.ParseInt(valueOf, 10, 64)
		if err != nil {
			return config.Terminal{}, errors.New("--server-port must be an integer(64)")
		}
		target.Config.Server.Port = int(serverPort)
	}
	valueOf, _ = cmd.Flags().GetString("server-mode")
	if valueOf != "" {
		target.Server.Mode = valueOf
	}
	return target, err
}
