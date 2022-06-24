package cmd

import (
	"github.com/liamg/dismember/internal/pkg/debug"
	"github.com/spf13/cobra"
)

var logger *debug.Logger

var rootCmd = &cobra.Command{
	Use:           "dismember",
	Short:         "",
	Long:          ``,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.SilenceUsage = true
		if flagDebug {
			logger = debug.New(cmd.ErrOrStderr())
		}
	},
}

var flagDebug bool

func Execute() error {
	rootCmd.PersistentFlags().BoolVarP(&flagDebug, "debug", "D", false, "Enable debug logging")
	return rootCmd.Execute()
}
