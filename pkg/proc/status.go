package proc

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// State represents the state of a process.
type State uint8

const (
	StateRunning     = 'R'
	StateSleeping    = 'S'
	StateDiskSleep   = 'D'
	StateStopped     = 'T'
	StateTracingStop = 't'
	StateZombie      = 'Z'
	StateDead        = 'X'
	StateUnknown     = 0
)

// String returns a string representation of the state.
func (s State) String() string {
	switch s {
	case StateRunning:
		return "R (running)"
	case StateSleeping:
		return "S (sleeping)"
	case StateDiskSleep:
		return "D (disk sleep)"
	case StateStopped:
		return "T (stopped)"
	case StateTracingStop:
		return "t (tracing stop)"
	case StateZombie:
		return "Z (zombie)"
	case StateDead:
		return "X (dead)"
	default:
		return "? (unknown)"
	}
}

// Status summarised data from /proc/[pid]/stat*
// fields with proc annotation are sourced from /proc/[pid]/status instead of /proc/[pid]/stat
type Status struct {
	Name                           string  `proc:"Name"` // Name of the command run by this process.  Strings longer than TASK_COMM_LEN (16) characters (including the terminating null byte) are silently truncated.
	State                          State   // Constant derived from StateDescription
	Parent                         Process // Parent process (0 if none)
	ProcessGroup                   int     // The process group ID
	Session                        int     // Session ID
	TTY                            Device  // The  controlling terminal of the process.  (The minor device number is contained in the combination of bits 31 to 20 and 7 to  0;  the  major  device number is in bits 15 to 8.)
	ForegroundTerminalProcessGroup int     // The  ID  of the foreground process group of the controlling terminal of the process.
	KernelFlags                    uint    // The kernel flags word of the process.  For bit meanings, see the  PF_*  defines  in  the Linux kernel source file include/linux/sched.h.  Details depend on the kernel version.
}

// State returns the state of the Process.
func (p *Process) State() State {
	status, err := p.Status()
	if err != nil {
		return StateUnknown
	}
	return status.State
}

// Status returns the status of the Process.
func (p *Process) Status() (*Status, error) {
	data, err := p.readFile("stat")
	if err != nil {
		return nil, err
	}
	status, err := parseStat(data)
	if err != nil {
		return nil, err
	}
	data, err = p.readFile("status")
	if err != nil {
		return nil, err
	}
	if err := parseStatus(data, status); err != nil {
		return nil, err
	}
	return status, nil
}

// IsAncestor returns true if the process is an ancestor of the given process, or if the process is the same as the given process.
func (p *Process) IsAncestor(other Process) bool {
	for other != NoProcess {
		if other == *p {
			return true
		}
		stat, err := other.Status()
		if err != nil {
			return false
		}
		other = stat.Parent
	}
	return false
}

func parseStat(data []byte) (*Status, error) {

	// prepend a blank entry, so we can use the indexes in `man proc`
	fields := append([]string{""}, strings.Fields(string(data))...)

	ppid, err := strconv.Atoi(fields[4])
	if err != nil {
		return nil, fmt.Errorf("invalid ppid '%s': %s", fields[4], err)
	}

	pgrp, err := strconv.Atoi(fields[5])
	if err != nil {
		return nil, fmt.Errorf("invalid pgrp '%s': %s", fields[5], err)
	}

	session, err := strconv.Atoi(fields[6])
	if err != nil {
		return nil, fmt.Errorf("invalid session '%s': %s", fields[6], err)
	}

	tty, err := strconv.Atoi(fields[7])
	if err != nil {
		return nil, fmt.Errorf("invalid tty '%s': %s", fields[7], err)
	}

	tpgid, err := strconv.Atoi(fields[8])
	if err != nil {
		return nil, fmt.Errorf("invalid tpgid '%s': %s", fields[8], err)
	}

	flags, err := strconv.Atoi(fields[9])
	if err != nil {
		return nil, fmt.Errorf("invalid flags '%s': %s", fields[9], err)
	}

	status := Status{
		Name:                           fields[2],
		State:                          State(fields[3][0]),
		Parent:                         Process(ppid),
		ProcessGroup:                   pgrp,
		Session:                        session,
		TTY:                            NewCharDeviceFromCombinedVersion(uint64(tty)),
		ForegroundTerminalProcessGroup: tpgid,
		KernelFlags:                    uint(flags),
	}

	return &status, nil
}

func parseStatus(data []byte, status *Status) error {

	values := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		val = strings.TrimSpace(val)
		values[key] = val
	}

	v := reflect.ValueOf(status)

	t := v.Elem().Type()
	for i := 0; i < t.NumField(); i++ {
		fv := t.Field(i)
		tags := strings.Split(fv.Tag.Get("proc"), ",")
		if len(tags) == 0 {
			continue
		}
		tagName := tags[0]
		if tagName == "-" {
			continue
		}
		value, ok := values[tagName]
		if !ok {
			continue
		}
		subject := v.Elem().Field(i)

		if !v.Elem().CanSet() {
			return fmt.Errorf("target is not settable")
		}

		switch subject.Kind() {
		case reflect.String:
			subject.SetString(value)
		case reflect.Uint64:
			u, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			subject.SetUint(uint64(u))
		default:
			return fmt.Errorf("decoding of kind %s is not supported", subject.Kind())
		}
	}
	return nil
}
