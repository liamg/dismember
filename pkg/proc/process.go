package proc

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type Process uint64

const (
	NoProcess Process = 0
)

func (p *Process) PID() uint64 {
	return uint64(*p)
}

func Self() Process {
	return Process(os.Getpid())
}

func (p *Process) Name() string {
	status, err := p.Status()
	if err != nil || status.Name == "" {
		return "unknown"
	}
	return status.Name
}

func (p *Process) String() string {
	return fmt.Sprintf("%d (%s)", p.PID(), p.Name())
}

var pidRegex = regexp.MustCompile(`^\d+$`)

func List(includeSelf bool) ([]Process, error) {

	self := os.Getpid()

	var results []Process
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !pidRegex.MatchString(entry.Name()) {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		if pid == self && !includeSelf {
			continue
		}
		results = append(results, Process(pid))
	}
	return results, nil
}

func (p *Process) readFile(path ...string) ([]byte, error) {
	final := filepath.Join(append([]string{"/proc", strconv.Itoa(int(p.PID()))}, path...)...)
	return os.ReadFile(final)
}

func (p *Process) openFile(path ...string) (*os.File, error) {
	final := filepath.Join(append([]string{"/proc", strconv.Itoa(int(p.PID()))}, path...)...)
	return os.Open(final)
}
