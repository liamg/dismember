package proc

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"syscall"
)

// Process represents a process. The underlying data is the PID.
type Process uint64

const (
	// NoProcess is a sentinel value for Process.
	NoProcess Process = 0
)

// PID returns the PID of the process.
func (p *Process) PID() uint64 {
	return uint64(*p)
}

// Self returns the Process of the currently running program
func Self() Process {
	return Process(os.Getpid())
}

// Ownership represents the ownership of a process.
type Ownership struct {
	UID uint32
	GID uint32
}

// Ownership returns the ownership of the process.
func (p *Process) Ownership() (*Ownership, error) {
	info, err := os.Stat(fmt.Sprintf("/proc/%d", *p))
	if err != nil {
		return nil, err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("stat syscall returned unexpected data: %#v", info.Sys())
	}
	return &Ownership{
		UID: stat.Uid,
		GID: stat.Gid,
	}, nil
}

// Name returns the name of the process.
func (p *Process) Name() string {
	status, err := p.Status()
	if err != nil || status.Name == "" {
		return "unknown"
	}
	return status.Name
}

// String returns the string representation of the process.
func (p *Process) String() string {
	return fmt.Sprintf("%d (%s)", p.PID(), p.Name())
}

var pidRegex = regexp.MustCompile(`^\d+$`)

// List returns a list of all processes available to the current user.
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
