package cmd

import (
	"fmt"
	"strings"

	"github.com/liamg/dismember/pkg/proc"
	"github.com/liamg/dismember/pkg/secrets"
	"github.com/spf13/cobra"
)

func init() {

	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "Search process memory for a set of predefined secret patterns",
		Long:  ``,
		RunE:  scanHandler,
		Args:  cobra.ExactArgs(0),
	}

	scanCmd.Flags().IntVarP(&flagPID, "pid", "p", 0, "PID of the process whose memory should be grepped. Omitting this option will grep the memory of all available processes on the system.")
	scanCmd.Flags().StringVarP(&flagProcessName, "process-name", "n", "", "Grep memory of all processes whose name contains this string.")
	scanCmd.Flags().IntVarP(&flagDumpRadius, "dump-radius", "r", 2, "The number of lines of memory to dump both above and below each match.")
	scanCmd.Flags().BoolVarP(&flagIncludeSelf, "self", "s", false, "Include results that are matched against the current process, or an ancestor of that process.")
	scanCmd.Flags().BoolVarP(&flagFast, "fast", "f", false, "Skip memory-mapped files in order to run faster.")
	rootCmd.AddCommand(scanCmd)
}

func scanHandler(cmd *cobra.Command, _ []string) error {

	var processes []proc.Process

	if flagPID == 0 {
		var err error
		processes, err = proc.List(false)
		if err != nil {
			return err
		}
	} else {
		processes = []proc.Process{proc.Process(flagPID)}
	}

	patterns := secrets.Patterns()

	stdErr := cmd.ErrOrStderr()
	_ = stdErr
	stdOut := cmd.OutOrStdout()

	var allResults []GrepResult

	for _, process := range processes {
		if flagProcessName != "" {
			status, err := process.Status()
			if err != nil {
				logger.Log("failed to determine status for process %s: %s", process, err)
				continue
			}
			if !strings.Contains(status.Name, flagProcessName) {
				continue
			}
		}
		if !flagIncludeSelf && process.IsAncestor(proc.Self()) {
			continue
		}
		for _, pattern := range patterns {
			results, err := grepProcessMemory(process, pattern)
			if err != nil {
				logger.Log("failed to search memory process %s: %s", process, err)
				continue
			}
			allResults = append(allResults, results...)
		}
	}

	if len(allResults) == 0 {
		_, _ = fmt.Fprintf(stdOut, "%sOperation Complete. No results found.%s\n\n", ansiRed, ansiReset)
	} else {
		for i, result := range allResults {
			_, _ = fmt.Fprint(stdOut, summariseResult(i+1, result))
		}
		_, _ = fmt.Fprintf(stdOut, "%sOperation Complete. %s%d%s%s results found.%s\n\n", ansiGreen, ansiBold, len(allResults), ansiReset, ansiGreen, ansiReset)
	}

	return nil
}
