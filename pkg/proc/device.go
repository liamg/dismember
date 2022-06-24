package proc

// Device represents a device in the /dev directory.
type Device struct {
	Major uint16
	Minor uint16
	Name  string
	Char  bool
}

// String returns a string representation of the device.
func (d Device) String() string {
	return d.Name
}
