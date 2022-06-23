package proc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Map(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Maps
		wantErr bool
	}{
		{
			input: `
55d38ffaf000-55d38ffb1000 r--p 00000000 103:03 16528031                  /usr/bin/cat
55d38ffb1000-55d38ffb5000 r-xp 00002000 103:03 16528031                  /usr/bin/cat
55d391c13000-55d391c34000 rw-p 00000000 00:00 0                          [heap]
7fb79f1ff000-7fb79f22e000 rw-p 00000000 00:00 0
ffffffffff600000-ffffffffff601000 --xp 00000000 00:00 0                  [vsyscall]
`,
			want: Maps{
				{
					Address: 0x55d38ffaf000,
					Size:    0x55d38ffb1000 - 0x55d38ffaf000,
					Permissions: MemPerms{
						Readable:   true,
						Writable:   false,
						Executable: false,
						Shared:     false,
					},
					Offset: 0,
					Device: (259 << 32) | 3,
					Inode:  16528031,
					Path:   "/usr/bin/cat",
				},
				{
					Address: 0x55d38ffb1000,
					Size:    0x55d38ffb5000 - 0x55d38ffb1000,
					Permissions: MemPerms{
						Readable:   true,
						Writable:   false,
						Executable: true,
						Shared:     false,
					},
					Offset: 0x2000,
					Device: (259 << 32) | 3,
					Inode:  16528031,
					Path:   "/usr/bin/cat",
				},
				{
					Address: 0x55d391c13000,
					Size:    0x55d391c34000 - 0x55d391c13000,
					Permissions: MemPerms{
						Readable:   true,
						Writable:   true,
						Executable: false,
						Shared:     false,
					},
					Offset: 0,
					Device: 0,
					Inode:  0,
					Path:   "[heap]",
				},
				{
					Address: 0x7fb79f1ff000,
					Size:    0x7fb79f22e000 - 0x7fb79f1ff000,
					Permissions: MemPerms{
						Readable:   true,
						Writable:   true,
						Executable: false,
						Shared:     false,
					},
					Offset: 0,
					Device: 0,
					Inode:  0,
					Path:   "",
				},
				{
					Address: 0xffffffffff600000,
					Size:    0xffffffffff601000 - 0xffffffffff600000,
					Permissions: MemPerms{
						Readable:   false,
						Writable:   false,
						Executable: true,
						Shared:     false,
					},
					Offset: 0,
					Device: 0,
					Inode:  0,
					Path:   "[vsyscall]",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseMaps([]byte(test.input))
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}
