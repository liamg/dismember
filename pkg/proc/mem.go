package proc

// ReadMemory reads the memory of the process for the given memory Map.
func (p *Process) ReadMemory(m Map, offset uint64, size uint64) ([]byte, error) {
	f, err := p.openFile("mem")
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	if _, err := f.Seek(int64(m.Address+offset), 0); err != nil {
		return nil, err
	}

	if size == 0 {
		size = m.Size
	}

	data := make([]byte, size)
	if _, err := f.Read(data); err != nil {
		return nil, err
	}

	return data, nil
}
