package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var logger *slog.Logger

type CommandConfig struct {
	Info string `yaml:"_info"`
	Cmd  string `yaml:"_cmd"`
}

var rootCmd = &cobra.Command{
	Use:     "na",
	Short:   "Dynamically builds CLI based on a YAML config",
	Version: "0.1.0",
}

func init() {
	logLevel := slog.LevelInfo
	if strings.TrimSpace(strings.ToLower(os.Getenv("SODIUM_LOG_LEVEL"))) == "debug" {
		logLevel = slog.LevelDebug
	}
	logOpts := &slog.HandlerOptions{
		Level: logLevel,
	}

	logger = slog.New(slog.NewTextHandler(os.Stdout, logOpts))

	// read config file
	var configPath string
	if envPath, ok := os.LookupEnv("SODIUM_CONFIG"); ok {
		configPath = envPath
	} else {
		configPath = filepath.Join(os.Getenv("HOME"), ".config", "sodium", ".narc.yaml")
		_, err := os.ReadFile(configPath)
		if err != nil {
			configPath = filepath.Join(os.Getenv("HOME"), ".config", "sodium", ".narc.yml")
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read the config file %v: %v\n", configPath, err)
		os.Exit(1)
	}

	aliases := make(map[string]interface{})
	if err := yaml.Unmarshal(data, &aliases); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal config file YAML: %v\n", err)
		os.Exit(1)
	}

	createCommands(aliases["aliases"].(map[interface{}]interface{}), rootCmd)
}

func createCommands(data map[interface{}]interface{}, parentCmd *cobra.Command) {
	for k, v := range data {
		key := k.(string)
		// handle minimal format command
		if cmd, ok := v.(string); ok {
			if strings.HasPrefix(key, "_") {
				continue
			}
			logger.Debug("parsed alias (minimal form)", "key", k, "cmd", cmd)
			cmd := &cobra.Command{
				Use:   key,
				Short: cmd,
				Run: func(c *cobra.Command, args []string) {
					// Extract the _cmd and execute it

					executeShellCmd(cmd, args)

				},
			}
			parentCmd.AddCommand(cmd)
			// handle full format command
		} else if subCmd, ok := v.(map[interface{}]interface{}); ok {
			subCmdInfo, infoExists := subCmd["_info"]
			cmdStr, cmdExists := subCmd["_cmd"].(string)

			if !infoExists { // add default _info based on _cmd
				subCmdInfo = ""
				if cmdExists {
					subCmdInfo = cmdStr
				}
			}

			if cmdExists {
				logger.Debug("parsed alias", "key", k, "info", subCmdInfo, "cmd", cmdStr)
			} else {
				logger.Debug("parsed group", "key", k, "info", subCmdInfo)
			}

			cmd := &cobra.Command{
				Use:   key,
				Short: subCmdInfo.(string),
				Run: func(c *cobra.Command, args []string) {
					// Extract the _cmd and execute it
					if cmdExists {
						executeShellCmd(cmdStr, args)
					}
				},
			}

			parentCmd.AddCommand(cmd)
			createCommands(subCmd, cmd)
		}
	}
}

func executeShellCmd(command string, args []string) {
	cmdWithArgs := strings.Split(command, " ")
	cmdWithArgs = append(cmdWithArgs, args...)

	logger.Debug("Executing command", "cmd", cmdWithArgs)

	cmd := exec.Command(cmdWithArgs[0], cmdWithArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Debug("Failed to execute the command", "err", err)
		exitErr, ok := err.(*exec.ExitError)
		if !ok {
			logger.Debug("Cannot extract original exit code, exit code will be -1", "err", err)
			os.Exit(-1)
		}
		os.Exit(exitErr.ExitCode()) // exit with original error code
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Debug("Failed to execute rootCmd", "err", err)
		os.Exit(1)
	}
}
