package proc

import "fmt"

// See https://www.kernel.org/doc/html/latest/admin-guide/devices.html#linux-allocated-devices-4-x-version

// NewCharDeviceFromCombinedVersion creates a Device from a combined version number.
func NewCharDeviceFromCombinedVersion(v uint64) Device {
	minor := ((v >> 16) & 0xff00) | (v & 0xff)
	major := (v >> 8) & 0xffff
	return NewCharDeviceFromVersion(uint16(major), uint16(minor))
}

// NewCharDeviceFromVersion creates a Device from a major and minor number.
func NewCharDeviceFromVersion(major, minor uint16) Device {
	return Device{
		Major: major,
		Minor: minor,
		Name:  lookupCharDeviceName(major, minor),
		Char:  true,
	}
}

func lookupCharDeviceName(major uint16, minor uint16) string {
	switch major {
	case 0:
		return "none"
	case 1:
		switch minor {
		case 1:
			return "/dev/mem"
		case 2:
			return "/dev/kmem"
		case 3:
			return "/dev/null"
		case 4:
			return "/dev/port"
		case 5:
			return "/dev/zero"
		case 6:
			return "/dev/core"
		case 7:
			return "/dev/full"
		case 8:
			return "/dev/random"
		case 9:
			return "/dev/urandom"
		case 10:
			return "/dev/aio"
		case 11:
			return "/dev/kmsg"
		case 12:
			return "/dev/oldmem"
		}
	case 2:
		switch minor {
		case 255:
			return "/dev/ptyef"
		default:
			return fmt.Sprintf("/dev/ptyp%d", minor)
		}
	case 3:
		switch minor {
		case 255:
			return "/dev/ttyef"
		default:
			return fmt.Sprintf("/dev/ttyp%d", minor)
		}
	case 4:
		switch {
		case minor >= 64:
			return fmt.Sprintf("/dev/ttyS%d", minor-64)
		default:
			return fmt.Sprintf("/dev/tty%d", minor)
		}
	case 5:
		switch minor {
		case 0:
			return "/dev/tty"
		case 1:
			return "/dev/console"
		case 2:
			return "/dev/ptmx"
		case 3:
			return "/dev/ttyprintk"
		default:
			if minor >= 64 {
				return fmt.Sprintf("/dev/cua%d", minor-64)
			}
		}
	case 7:
		switch {
		case minor == 0:
			return "/dev/vcs"
		case minor < 64:
			return fmt.Sprintf("/dev/vcs%d", minor)
		case minor == 64:
			return "/dev/vcsu"
		case minor > 64:
			return fmt.Sprintf("/dev/vcsu%d", minor-64)
		}
	case 136, 137, 138, 139, 140, 141, 142, 143:
		return fmt.Sprintf("/dev/pts/%d", minor)
	}
	return fmt.Sprintf("unknown (%d/%d)", major, minor)
}
