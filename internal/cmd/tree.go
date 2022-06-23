package cmd

import (
	"fmt"
	"github.com/liamg/dismember/pkg/proc"
	"github.com/spf13/cobra"
	"io"
)

func init() {
	treeCmd := &cobra.Command{
		Use:   "tree",
		Short: "Show a tree diagram of a process and all children (defaults to PID 1).",
		Long:  ``,
		RunE:  treeHandler,
		Args:  cobra.ExactArgs(0),
	}
	treeCmd.Flags().IntVarP(&flagPID, "pid", "p", 1, "PID of the process to analyse")
	rootCmd.AddCommand(treeCmd)

}

type procWithStatus struct {
	process proc.Process
	status  proc.Status
}

func treeHandler(cmd *cobra.Command, _ []string) error {

	processes, err := proc.List(true)
	if err != nil {
		return err
	}

	root := proc.Process(flagPID)
	status, err := root.Status()
	if err != nil {
		return err
	}
	rootWithStatus := procWithStatus{
		process: root,
		status:  *status,
	}

	var all []procWithStatus
	for _, process := range processes {
		status, err := process.Status()
		if err != nil {
			continue
		}
		all = append(all, procWithStatus{
			process: process,
			status:  *status,
		})
	}

	drawBranch(cmd.OutOrStdout(), rootWithStatus, nil, all)
	return nil
}

func drawBranch(w io.Writer, parent procWithStatus, lasts []bool, all []procWithStatus) {

	var children []procWithStatus
	for _, process := range all {
		if process.status.Parent != parent.process {
			continue
		}
		children = append(children, process)
	}

	var done bool
	_, _ = fmt.Fprint(w, ansiDim)
	if len(lasts) > 1 {
		for _, last := range lasts[1:] {
			if !last {
				_, _ = fmt.Fprint(w, " │ ")
			} else if !done {
				done = true
				_, _ = fmt.Fprint(w, " │ ")
			} else {
				_, _ = fmt.Fprint(w, "   ")
			}
		}
	}

	if len(lasts) > 0 {
		symbol := '├'
		if lasts[len(lasts)-1] {
			symbol = '└'
		}
		_, _ = fmt.Fprintf(w, " %c─ ", symbol)
	}

	_, _ = fmt.Fprint(w, ansiReset)
	_, _ = fmt.Fprintf(w, "%s %s(%s%d%s)%s\n", parent.status.Name, ansiDim, ansiReset, parent.process, ansiDim, ansiReset)

	for i, child := range children {
		drawBranch(w, child, append(lasts, i == len(children)-1), all)
	}
}
