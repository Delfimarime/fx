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

func RunSpf13CobraHttpServer(appName, description string, f func(config.Terminal) []fx.Option, fxConfig ...fx.Option) {
	rootCmd := NewSpf13CobraHttpServer(appName, description, f, fxConfig...)
	if err := rootCmd.Execute(); err != nil {
		zap.L().Error("Stopping server...", zap.Error(err))
		os.Exit(1)
	}
}

func NewSpf13CobraHttpServer(appName, description string, f func(config.Terminal) []fx.Option, fxConfig ...fx.Option) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   appName,
		Short: description,
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
			cmd.Println(fmt.Sprintf("%s\nversion:%s", appName, version))
			return nil
		},
	}
	startupCmd := &cobra.Command{
		Use:   "start",
		Short: fmt.Sprintf("Starts the HTTP Server that exposes %s as an API.", appName),
		Args:  cobra.MaximumNArgs(1),
		Run:   NewSpf13CobraHttpServerCommand(appName, f, fxConfig...),
	}
	// Add flags to the startup command
	startupCmd.PersistentFlags().String("log-level", "INFO", "Set the logging level (e.g., DEBUG, INFO)")
	startupCmd.PersistentFlags().Int("server-port", 8080, "Port on which the server runs")
	startupCmd.PersistentFlags().String("server-mode", "release", "The mode the server is running")
	rootCmd.AddCommand(startupCmd)
	return rootCmd
}

func NewSpf13CobraHttpServerCommand(appName string, f func(config.Terminal) []fx.Option, fxConfig ...fx.Option) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		configuration, err := getSpf13CobraTerminalConfiguration(cmd, args, appName)
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
		opts := make([]fx.Option, 0)
		if f != nil {
			opts = f(configuration)
		}
		if fxConfig != nil {
			opts = append(opts, fxConfig...)
		}
		fx.New(append(opts, fx.Provide(func() config.Terminal {
			return configuration
		}, func(c config.Terminal) config.Config {
			return c.Config
		}))...).Run()
	}
}

func getSpf13CobraTerminalConfiguration(cmd *cobra.Command, args []string, appName string) (config.Terminal, error) {
	configURI := ""
	if len(args) > 0 {
		configURI = args[0]
	}
	if configURI == "" {
		value := os.Getenv(strings.ToUpper(fmt.Sprintf("%s_HOME", appName)))
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
