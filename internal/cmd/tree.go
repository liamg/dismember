package cmd

import (
	"fmt"
	"github.com/liamg/dismember/pkg/proc"
	"github.com/spf13/cobra"
	"io"
	"os"
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

	uid := os.Getuid()
	drawBranch(cmd.OutOrStdout(), rootWithStatus, "", true, all, uid)
	return nil
}

func drawBranch(w io.Writer, parent procWithStatus, prefix string, last bool, all []procWithStatus, uid int) {

	var children []procWithStatus
	for _, process := range all {
		if process.status.Parent != parent.process {
			continue
		}
		children = append(children, process)
	}

	_, _ = fmt.Print(ansiDim + prefix)
	if prefix != "" {
		symbol := '├'
		if last {
			symbol = '└'
		}
		_, _ = fmt.Fprintf(w, " %c─ ", symbol)
	}

	_, _ = fmt.Fprint(w, ansiReset)
	//owner, err := parent.process.Ownership()
	//if err != nil {
	//
	//}
	_, _ = fmt.Fprintf(w, "%s %s(%s%d%s)%s\n", parent.status.Name, ansiDim, ansiReset, parent.process, ansiDim, ansiReset)
	_, _ = fmt.Fprint(w, ansiReset)

	if last {
		prefix += "   "
	} else {
		prefix += " │ "
	}

	for i, child := range children {
		drawBranch(w, child, prefix, i == len(children)-1, all)
	}
}
