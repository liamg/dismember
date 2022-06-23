package proc

import (
	"fmt"
	"strconv"
	"strings"
)

// see https://man7.org/linux/man-pages/man5/proc.5.html
//
// Example:
// address           perms offset  dev   inode   pathname
// 08048000-08056000 r-xp 00000000 03:0c 64593   /usr/sbin/gpm

type Maps []Map

type Map struct {
	Address     uint64
	Size        uint64
	Permissions MemPerms
	Offset      uint64
	Device      uint64 // see man makedev
	Inode       uint64
	Path        string
}

type MemPerms struct {
	Readable   bool
	Writable   bool
	Executable bool
	Shared     bool
}

func (p *Process) Maps() (Maps, error) {
	data, err := p.readFile("maps")
	if err != nil {
		return nil, err
	}
	return parseMaps(data)
}

func parseMaps(data []byte) (Maps, error) {
	var maps Maps
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		if len(fields) == 5 {
			fields = append(fields, "")
		}
		var m Map
		start, end, err := parseAddressRange(fields[0])
		if err != nil {
			return nil, err
		}
		m.Address = start
		m.Size = end - start
		m.Permissions, err = parsePermissions(fields[1])
		if err != nil {
			return nil, err
		}
		m.Offset, err = parseUint64Hex(fields[2])
		if err != nil {
			return nil, err
		}
		m.Device, err = parseMkdev(fields[3])
		if err != nil {
			return nil, err
		}
		m.Inode, err = parseUint64Dec(fields[4])
		if err != nil {
			return nil, err
		}
		m.Path = fields[5]
		maps = append(maps, m)
	}
	return maps, nil
}

func parsePermissions(s string) (MemPerms, error) {
	var perms MemPerms
	if len(s) != 4 {
		return perms, fmt.Errorf("invalid permissions: %s", s)
	}
	perms.Readable = s[0] == 'r'
	perms.Writable = s[1] == 'w'
	perms.Executable = s[2] == 'x'
	perms.Shared = s[3] == 's'
	return perms, nil
}

func parseAddressRange(input string) (uint64, uint64, error) {
	fields := strings.Split(input, "-")
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("invalid address range: %s", input)
	}
	start, err := parseUint64Hex(fields[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start address: %s: %w", fields[0], err)
	}
	end, err := parseUint64Hex(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start address: %s: %w", fields[0], err)
	}
	return start, end, nil
}

func parseUint64Hex(input string) (uint64, error) {
	return strconv.ParseUint(input, 16, 64)
}

func parseUint64Dec(input string) (uint64, error) {
	return strconv.ParseUint(input, 10, 64)
}

func parseMkdev(input string) (uint64, error) {
	major, minor, ok := strings.Cut(input, ":")
	if !ok {
		return parseUint64Hex(input)
	}
	ma, err := parseUint64Hex(major)
	if err != nil {
		return 0, err
	}
	mi, err := parseUint64Hex(minor)
	if err != nil {
		return 0, err
	}
	return (ma << 32) | mi, nil
}
