package cmd

import (
	"fmt"
	"strconv"

	"github.com/liamg/dismember/pkg/proc"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "files [pid]",
		Short: "Show a list of files being accessed by a process",
		RunE:  filesHandler,
		Args:  cobra.ExactArgs(1),
	})
}

func filesHandler(cmd *cobra.Command, args []string) error {

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid pid specified: '%s': %w", args[0], err)
	}

	process := proc.Process(pid)

	files, err := process.Files()
	if err != nil {
		return fmt.Errorf("failed to read accessed files for process %d: %w\n", process.PID(), err)
	}

	w := cmd.OutOrStdout()
	for _, file := range files {
		_, _ = fmt.Fprintln(w, file)
	}
	return nil
}
