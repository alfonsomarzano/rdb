package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	repoPath string
	jsonOutput bool
	noColor bool
	trace bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rdb",
	Short: "Resource Database CLI tool",
	Long: `RDB is a Windows-first CLI tool for managing a local "Resource Database" 
— a directory tree of typed assets — with git-like workflows.

Features:
- Typed assets stored under folders named by ID (e.g., 1030002/, 1000624/, 1010042/)
- Git-like workflows (init, status, add, commit, branch, merge)
- Content integrity with SHA-256 for objects
- Human-readable directory layout
- Portable .rdbdata packages`,
	Version: "1.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .rdb/config.json)")
	rootCmd.PersistentFlags().StringVar(&repoPath, "repo", "", "operate on repo at path (default: cwd)")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "JSON output for scripts")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable color output")
	rootCmd.PersistentFlags().BoolVar(&trace, "trace", false, "verbose internal logging")

	// Bind flags to viper
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("repo", rootCmd.PersistentFlags().Lookup("repo"))
	viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
	viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))
	viper.BindPFlag("trace", rootCmd.PersistentFlags().Lookup("trace"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in .rdb directory with name "config" (without extension).
		viper.AddConfigPath(".rdb")
		viper.SetConfigType("json")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if trace {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
} 