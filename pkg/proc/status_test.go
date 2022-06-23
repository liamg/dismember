package proc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Status(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Status
		wantErr bool
	}{
		{
			input: `
Name:	cat
Umask:	0022
State:	R (running)
Tgid:	131533
Ngid:	0
Pid:	131533
PPid:	118575
TracerPid:	0
Uid:	0	0	0	0
Gid:	0	0	0	0
FDSize:	256
Groups:	0
NStgid:	131533
NSpid:	131533
NSpgid:	131533
NSsid:	118559
VmPeak:	    5788 kB
VmSize:	    5788 kB
VmLck:	       0 kB
VmPin:	       0 kB
VmHWM:	     976 kB
VmRSS:	     976 kB
RssAnon:	      88 kB
RssFile:	     888 kB
RssShmem:	       0 kB
VmData:	     360 kB
VmStk:	     132 kB
VmExe:	      16 kB
VmLib:	    1668 kB
VmPTE:	      48 kB
VmSwap:	       0 kB
HugetlbPages:	       0 kB
CoreDumping:	0
THP_enabled:	1
Threads:	1
SigQ:	1/127176
SigPnd:	0000000000000000
ShdPnd:	0000000000000000
SigBlk:	0000000000000000
SigIgn:	0000000000000000
SigCgt:	0000000000000000
CapInh:	0000000000000000
CapPrm:	000001ffffffffff
CapEff:	000001ffffffffff
CapBnd:	000001ffffffffff
CapAmb:	0000000000000000
NoNewPrivs:	0
Seccomp:	0
Seccomp_filters:	0
Speculation_Store_Bypass:	thread vulnerable
SpeculationIndirectBranch:	conditional enabled
Cpus_allowed:	ff
Cpus_allowed_list:	0-7
Mems_allowed:	00000001
Mems_allowed_list:	0
voluntary_ctxt_switches:	1
nonvoluntary_ctxt_switches:	0
`,
			want: Status{
				Name:   "cat",
				Parent: Process(118575),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseStatus([]byte(test.input))
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, *got)
		})
	}
}
