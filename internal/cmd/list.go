package cmd

import (
	"fmt"

	"github.com/liamg/dismember/pkg/proc"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all processes currently available on the system",
		Long:  ``,
		RunE:  listHandler,
		Args:  cobra.ExactArgs(0),
	})
}

func listHandler(cmd *cobra.Command, _ []string) error {

	processes, err := proc.List(true)
	if err != nil {
		return err
	}

	stdErr := cmd.ErrOrStderr()
	stdOut := cmd.OutOrStdout()

	for _, process := range processes {
		status, err := process.Status()
		if err != nil {
			_, _ = fmt.Fprintf(stdErr, "failed to read status for process %d: %s\n", process.PID(), err)
			continue
		}
		_, _ = fmt.Fprintf(stdOut, "% -10d %s\n", process.PID(), status.Name)
	}

	return nil
}
