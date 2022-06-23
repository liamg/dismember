package cmd

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/liamg/dismember/pkg/proc"
	"github.com/spf13/cobra"
)

var flagPID int
var flagProcessName string
var flagIncludeSelf bool
var flagDumpRadius int
var flagFast bool

func init() {

	grepCmd := &cobra.Command{
		Use:   "grep [keyword]",
		Short: "Search process memory for a given string or regex",
		Long:  ``,
		RunE:  grepHandler,
		Args:  cobra.ExactArgs(1),
	}

	grepCmd.Flags().IntVarP(&flagPID, "pid", "p", 0, "PID of the process whose memory should be grepped. Omitting this option will grep the memory of all available processes on the system.")
	grepCmd.Flags().StringVarP(&flagProcessName, "pname", "n", "", "Grep memory of all processes whose name contains this string.")
	grepCmd.Flags().IntVarP(&flagDumpRadius, "dump-radius", "r", 2, "The number of lines of memory to dump both above and below each match.")
	grepCmd.Flags().BoolVarP(&flagIncludeSelf, "self", "s", false, "Include results that are matched against the current process, or an ancestor of that process.")
	grepCmd.Flags().BoolVarP(&flagFast, "fast", "f", false, "Skip mapped files in order to run faster.")
	rootCmd.AddCommand(grepCmd)
}

func grepHandler(cmd *cobra.Command, args []string) error {

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

	regex, err := regexp.Compile(args[0])
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}

	stdErr := cmd.ErrOrStderr()
	_ = stdErr
	stdOut := cmd.OutOrStdout()

	var total int
	for _, process := range processes {
		if flagProcessName != "" {
			status, err := process.Status()
			if err != nil {
				continue
			}
			if !strings.Contains(status.Name, flagProcessName) {
				continue
			}
		}
		if !flagIncludeSelf && process.IsInHierarchyOf(proc.Self()) {
			continue
		}
		results, err := grepProcessMemory(process, regex)
		if err != nil {
			// TODO: add to debug log
			//_, _ = fmt.Fprintf(stdErr, "failed to access memory for process %d: %s\n", process.PID(), err)
			continue
		}
		for i, result := range results {
			_, _ = fmt.Fprint(stdOut, summariseResult(total+i+1, result))
		}
		total += len(results)
	}
	if total == 0 {
		_, _ = fmt.Fprintf(stdOut, "%sOperation Complete. No results found.%s\n\n", ansiRed, ansiReset)
	} else {
		_, _ = fmt.Fprintf(stdOut, "%sOperation Complete. %s%d%s%s results found.%s\n\n", ansiGreen, ansiBold, total, ansiReset, ansiGreen, ansiReset)
	}

	return nil
}

type GrepResult struct {
	Pattern *regexp.Regexp
	Process proc.Process
	Map     proc.Map
	Address uint64
	Match   []byte
}

const (
	ansiReset     = "\x1b[0m"
	ansiBold      = "\x1b[1m"
	ansiDim       = "\x1b[2m"
	ansiItalic    = "\x1b[3m"
	ansiUnderline = "\x1b[4m"
	ansiRed       = "\x1b[31m"
	ansiGreen     = "\x1b[32m"
)

func summariseResult(number int, g GrepResult) string {

	buffer := bytes.NewBuffer(nil)

	_, _ = fmt.Fprintf(buffer, " %sMatch #%d%s\n\n", ansiUnderline, number, ansiReset)
	_, _ = fmt.Fprintf(buffer, "  %sMatched%s   %s\n", ansiBold, ansiReset, string(g.Match))
	_, _ = fmt.Fprintf(buffer, "  %sPattern%s   %s\n", ansiBold, ansiReset, g.Pattern.String())
	_, _ = fmt.Fprintf(buffer, "  %sProcess%s   %s\n", ansiBold, ansiReset, g.Process.String())
	_, _ = fmt.Fprintf(buffer, "  %sAddress%s   0x%x %s\n\n", ansiBold, ansiReset, g.Address, g.Map.Path)
	_, _ = fmt.Fprintf(buffer, "  %sMemory Dump%s\n\n%s\n\n", ansiBold, ansiReset, hexDump(g))

	return buffer.String()
}

func hexDump(g GrepResult) string {

	buffer := bytes.NewBuffer(nil)

	offset := g.Address - g.Map.Address

	linesEitherSide := uint64(flagDumpRadius)
	if linesEitherSide < 0 {
		linesEitherSide = 0
	}

	start := ((offset / 16) * 16) - (16 * linesEitherSide)
	size := (((uint64(len(g.Match)) + (16 * (2 * linesEitherSide))) / 16) + 1) * 16

	data, err := g.Process.ReadMemory(g.Map, start, size)
	if err != nil {
		return fmt.Sprintf("    dump not available: %s", err)
	}

	literalStartAddr := g.Map.Address + start

	_, _ = fmt.Fprintf(buffer, "                    %s", ansiDim)
	for i := 0; i < 0x10; i++ {
		_, _ = fmt.Fprintf(buffer, "%02X ", i)
	}
	_, _ = fmt.Fprintln(buffer, ansiReset)

	var ascii string
	for index, b := range data {

		localIndex := uint64(index) + start
		inSecret := localIndex >= offset && localIndex < offset+uint64(len(g.Match))

		if index%16 == 0 && index > 0 {
			_, _ = fmt.Fprintf(buffer, "  %s\n", ascii)
			ascii = ""
		}
		if index%16 == 0 {
			_, _ = fmt.Fprintf(buffer, "  %s%016x%s  ", ansiDim, literalStartAddr+uint64(index/16), ansiReset)
		}

		if inSecret {
			_, _ = fmt.Fprintf(buffer, "%s%s", ansiBold, ansiRed)
		}
		_, _ = fmt.Fprintf(buffer, "%02x%s ", b, ansiReset)
		ascii += asciify(b, inSecret)
	}
	if ascii != "" {
		_, _ = fmt.Fprintf(buffer, "  %s\n", ascii)
	}

	return buffer.String()
}

func asciify(b byte, hl bool) string {
	if b < ' ' || b > '~' {
		b = '.'
	}
	if !hl {
		return fmt.Sprintf("%s%c%s", ansiDim, b, ansiReset)
	}
	return fmt.Sprintf("%s%s%c%s", ansiBold, ansiRed, b, ansiReset)
}

func grepProcessMemory(p proc.Process, regex *regexp.Regexp) ([]GrepResult, error) {
	var results []GrepResult
	maps, err := p.Maps()
	if err != nil {
		return nil, err
	}
	for _, map_ := range maps {
		if !map_.Permissions.Readable {
			continue
		}
		memory, err := p.ReadMemory(map_, 0, 0)
		if err != nil {
			continue
		}
		for _, matches := range regex.FindAllIndex(memory, -1) {
			results = append(results, GrepResult{
				Process: p,
				Map:     map_,
				Address: map_.Address + uint64(matches[0]),
				Match:   shrinkMatch(memory[matches[0]:matches[1]]),
				Pattern: regex,
			})
		}
	}
	return results, nil
}

func shrinkMatch(match []byte) []byte {
	return bytes.Split(match, []byte{0x00})[0]
}
