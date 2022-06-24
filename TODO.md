# TODO

- [x] Add lots more info to status (ownership etc.)
- [x] Add command to find pid by name
- [x] Add suspend/resume commands (`syscall.Kill(process.PID(), syscall.SIGINT)`)
- [x] Add kill command (with --children flag)
- [x] Add kernel command to show kernel info (/proc/cmdline, /proc/sys/kernel/*)
- [x] Add files command to show list of files open by a process (/proc/123/fd/*)
- [x] Add godoc everywhere
- [x] Add scan command with built-in patterns
- [ ] Add debug logging