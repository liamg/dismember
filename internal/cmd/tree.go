package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/liamg/dismember/pkg/proc"
	"github.com/spf13/cobra"
)

var flagTreePID int

func init() {
	treeCmd := &cobra.Command{
		Use:   "tree",
		Short: "Show a tree diagram of a process and all children (defaults to PID 1).",
		Long:  ``,
		RunE:  treeHandler,
		Args:  cobra.ExactArgs(0),
	}
	treeCmd.Flags().IntVarP(&flagTreePID, "pid", "p", 1, "PID of the process to analyse")
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

	root := proc.Process(flagTreePID)
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
	drawBranch(cmd.OutOrStdout(), rootWithStatus, "", true, true, all, uid)
	return nil
}

func drawBranch(w io.Writer, parent procWithStatus, prefix string, first bool, last bool, all []procWithStatus, uid int) {

	var children []procWithStatus
	for _, process := range all {
		if process.status.Parent != parent.process {
			continue
		}
		children = append(children, process)
	}

	_, _ = fmt.Print(ansiDim + prefix)
	if !first {
		symbol := '├'
		if last {
			symbol = '└'
		}
		_, _ = fmt.Fprintf(w, " %c─ ", symbol)
	}

	_, _ = fmt.Fprint(w, ansiReset)
	ownerName := "uid=?"
	if owner, err := parent.process.Ownership(); err == nil {
		if int(owner.UID) == uid {
			ownerName = ""
		} else {
			ownerName = fmt.Sprintf("uid=%d", owner.UID)
		}
	}
	_, _ = fmt.Fprintf(w, "%s %s(%s%d%s)%s %s\n", parent.status.Name, ansiDim, ansiReset, parent.process, ansiDim, ansiReset, ownerName)
	_, _ = fmt.Fprint(w, ansiReset)

	if !first {
		if last {
			prefix += "   "
		} else {
			prefix += " │ "
		}
	}

	for i, child := range children {
		drawBranch(w, child, prefix, false, i == len(children)-1, all, uid)
	}
}
