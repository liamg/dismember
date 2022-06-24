package cmd

import (
	"fmt"

	"github.com/liamg/dismember/pkg/proc"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "find [name]",
		Short: "Find a PID given a process name. If multiple processes match, the first one is returned.",
		Long:  ``,
		RunE:  findHandler,
		Args:  cobra.ExactArgs(1),
	})
}

func findHandler(cmd *cobra.Command, args []string) error {

	processes, err := proc.List(true)
	if err != nil {
		return err
	}

	w := cmd.OutOrStdout()

	for _, process := range processes {
		status, err := process.Status()
		if err != nil {
			logger.Log("failed to read status for process %d: %s\n", process.PID(), err)
			continue
		}
		if status.Name == args[0] {
			_, _ = fmt.Fprintf(w, "%d\n", process.PID())
			return nil
		}
	}

	return fmt.Errorf("no process found with name '%s'", args[0])
}
