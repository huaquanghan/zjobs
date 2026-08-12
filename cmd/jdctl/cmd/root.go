package cmd

import (
	"github.com/spf13/cobra"
)

// version is injected at build time via -ldflags; "dev" is the local default.
var version = "dev"

// rootCmd is the stable entry surface for jdctl. Command names and primary
// flags are part of the locked MVP contract; do not rename without a plan change.
var rootCmd = &cobra.Command{
	Use:     "jdctl",
	Short:   "Local-first JD analyzer: profile + CV + JD in, fit report out",
	Version: version,
}

func init() {
	rootCmd.PersistentFlags().String("config", "", "path to YAML config (defaults apply when empty)")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
