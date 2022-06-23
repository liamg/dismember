package cmd

import (
	"fmt"
	"io"
	"strconv"

	"github.com/liamg/dismember/pkg/proc"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "status [pid]",
		Short: "Show information about the status of a process",
		Long:  ``,
		RunE:  statusHandler,
		Args:  cobra.ExactArgs(1),
	})
}

func statusHandler(cmd *cobra.Command, args []string) error {

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid pid specified: '%s': %w", args[0], err)
	}

	process := proc.Process(pid)

	status, err := process.Status()
	if err != nil {
		return fmt.Errorf("failed to read status for process %d: %w\n", process.PID(), err)
	}

	stdOut := cmd.OutOrStdout()

	printKeyVal(stdOut, "PID", strconv.Itoa(int(process.PID())))
	printKeyVal(stdOut, "Name", status.Name)
	if status.Parent != 0 {
		printKeyVal(stdOut, "Parent", status.Parent.String())
	} else {
		printKeyVal(stdOut, "Parent", "-")
	}

	return nil
}

func printKeyVal(w io.Writer, key string, value string) {
	_, _ = fmt.Fprintf(w, "%-20s %s\n", key, value)
}
