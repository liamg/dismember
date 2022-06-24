package cmd

import (
	"fmt"
	"strconv"
	"syscall"

	"github.com/liamg/dismember/pkg/proc"
	"github.com/spf13/cobra"
)

var flagKillChildren bool

func init() {
	killCmd := &cobra.Command{
		Use:   "kill [pid]",
		Short: "Kill a process using SIGKILL",
		RunE:  killHandler,
		Args:  cobra.ExactArgs(1),
	}
	killCmd.Flags().BoolVarP(&flagKillChildren, "children", "c", false, "Kill all children of the specified process (leaving the process itself alive)")
	rootCmd.AddCommand(killCmd)
}

func killHandler(cmd *cobra.Command, args []string) error {

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid pid specified: '%s': %w", args[0], err)
	}

	process := proc.Process(pid)

	if !flagKillChildren {
		if err := syscall.Kill(int(process.PID()), syscall.SIGKILL); err != nil {
			return fmt.Errorf("failed to kill process %d: %w\n", process.PID(), err)
		}

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Process %s killed.\n", process.String())
		return nil
	}

	processes, err := proc.List(true)
	if err != nil {
		return fmt.Errorf("failed to list children for process: %w", err)
	}

	for _, candidate := range processes {
		status, err := candidate.Status()
		if err != nil {
			logger.Log("failed to determine status for process %s: %s", candidate, err)
			continue
		}
		if status.Parent == process {
			if err := syscall.Kill(int(candidate.PID()), syscall.SIGKILL); err != nil {
				return fmt.Errorf("failed to kill child process %d: %w\n", candidate.PID(), err)
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Child process %s killed.\n", candidate.String())
		}
	}
	return nil
}
