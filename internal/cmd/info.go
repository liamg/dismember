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
		Use:   "info [pid]",
		Short: "Show information about a process",
		RunE:  infoHandler,
		Args:  cobra.ExactArgs(1),
	})
}

func infoHandler(cmd *cobra.Command, args []string) error {

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid pid specified: '%s': %w", args[0], err)
	}

	process := proc.Process(pid)

	status, err := process.Status()
	if err != nil {
		return fmt.Errorf("failed to read status for process %d: %w\n", process.PID(), err)
	}

	owner, err := process.Ownership()
	if err != nil {
		return fmt.Errorf("failed to read ownership for process %d: %w\n", process.PID(), err)
	}

	stdOut := cmd.OutOrStdout()

	printKeyValUint64Decimal(stdOut, "PID", process.PID())
	printKeyValString(stdOut, "Name", status.Name)
	printKeyValString(stdOut, "State", status.State.String())
	if status.Parent != 0 {
		printKeyValString(stdOut, "Parent", status.Parent.String())
	} else {
		printKeyValString(stdOut, "Parent", "-")
	}
	printKeyValUint64Decimal(stdOut, "Process Group", uint64(status.ProcessGroup))
	printKeyValUint64Decimal(stdOut, "Session", uint64(status.Session))
	printKeyValString(stdOut, "TTY", status.TTY.String())
	printKeyValUint64Decimal(stdOut, "Terminal Process Group", uint64(status.ForegroundTerminalProcessGroup))
	printKeyValUint64Hex(stdOut, "Kernel Flags", uint64(status.KernelFlags))
	printKeyValUint64Decimal(stdOut, "Owner UID", uint64(owner.UID))
	printKeyValUint64Decimal(stdOut, "Owner GID", uint64(owner.GID))
	return nil
}

func printKeyValString(w io.Writer, key string, value string) {
	_, _ = fmt.Fprintf(w, "%-24s %s\n", key, value)
}

func printKeyValUint64Decimal(w io.Writer, key string, value uint64) {
	_, _ = fmt.Fprintf(w, "%-24s %d\n", key, value)
}

func printKeyValUint64Hex(w io.Writer, key string, value uint64) {
	_, _ = fmt.Fprintf(w, "%-24s 0x%x\n", key, value)
}
