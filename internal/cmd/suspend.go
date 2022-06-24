package cmd

import (
	"fmt"
	"strconv"
	"syscall"

	"github.com/liamg/dismember/pkg/proc"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "suspend [pid]",
		Short: "Suspend a process using SIGSTOP (use 'dismember resume' to leave suspension)",
		RunE:  suspendHandler,
		Args:  cobra.ExactArgs(1),
	})
}

func suspendHandler(cmd *cobra.Command, args []string) error {

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid pid specified: '%s': %w", args[0], err)
	}

	process := proc.Process(pid)

	status, err := process.Status()
	if err != nil {
		return fmt.Errorf("failed to read status for process %d: %w\n", process.PID(), err)
	}

	if status.State == proc.StateStopped {
		return fmt.Errorf("process %d is already stopped", process.PID())
	}

	if err := syscall.Kill(int(process.PID()), syscall.SIGSTOP); err != nil {
		return fmt.Errorf("failed to suspend process %d: %w\n", process.PID(), err)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Process %s suspended. Use 'dismember resume %d' to resume,\n", process.String(), process)
	return nil
}
