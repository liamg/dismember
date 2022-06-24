package proc

import (
	"fmt"
	"os"
	"path/filepath"
)

// Files discovers a list of all files being accessed the Process.
func (p *Process) Files() ([]string, error) {
	base := fmt.Sprintf("/proc/%d/fd", p.PID())
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	var matches []string
	for _, entry := range entries {
		link, err := os.Readlink(filepath.Join(base, entry.Name()))
		if err != nil {
			continue
		}
		matches = append(matches, link)
	}
	return matches, nil
}
