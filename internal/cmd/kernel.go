package cmd

import (
	"fmt"

	"github.com/liamg/dismember/pkg/proc"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "kernel",
		Short: "Show information about the kernel",
		RunE:  kernelHandler,
		Args:  cobra.ExactArgs(0),
	})
}

func kernelHandler(cmd *cobra.Command, args []string) error {
	info := proc.GetKernelInfo()
	stdOut := cmd.OutOrStdout()
	printKeyValString(stdOut, "Type", info.OSType)
	printKeyValString(stdOut, "Release", info.OSRelease)
	printKeyValString(stdOut, "Boot Args", info.BootArgs)
	_, _ = fmt.Fprintf(stdOut, "\n%s\n", info.FullVersion)
	return nil
}
