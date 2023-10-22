package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type CommandConfig struct {
	Info string `yaml:"_info"`
	Cmd  string `yaml:"_cmd"`
}

var rootCmd = &cobra.Command{
	Use:   "na",
	Short: "Dynamically built CLI based on a YAML config",
}

func init() {
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "sodium", ".narc.yaml")
	if envPath, ok := os.LookupEnv("SODIUM_CONFIG"); ok {
		configPath = envPath
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

		// If it's a CommandConfig, add the actual command
		if subCmd, ok := v.(map[interface{}]interface{}); ok {
			subCmdInfo, infoExists := subCmd["_info"]
			if !infoExists { // add default _info based on _cmd
				subCmdInfo = ""
				if cmdStr, exists := subCmd["_cmd"].(string); exists {
					subCmdInfo = cmdStr
				}
			}
			cmd := &cobra.Command{
				Use:   key,
				Short: subCmdInfo.(string),
				Run: func(c *cobra.Command, args []string) {
					// Extract the _cmd and execute it
					if cmdStr, exists := subCmd["_cmd"].(string); exists {
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

	cmd := exec.Command(cmdWithArgs[0], cmdWithArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute the command: %v\n", err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute rootCmd: %v\n", err)
		os.Exit(1)
	}
}
