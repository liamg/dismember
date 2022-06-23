# dismember

Dismember is a command-line tool for Linux used to grep for patterns across the entire memory used by a process (or processes).

![A gif showing dismember finding a password from a Slack message](demo.gif)

Dismember can be used to search memory of all processes it has access to, so running it as root is the most effective method.

Commands are also included to list processes, explore process status and related information, draw process trees, and more...

## Installation

Grab a binary from the [latest release](https://github.com/liamg/dismember/releases/latest) and add it to your path.

## Examples

### Search for a pattern in a process by PID
```bash
# search memory owned by process 1234
dismember grep -p 1234 'the password is .*'
```

### Search for a pattern in a process by name
```bash
# search memory owned by processes named "nginx" for a login form submission
dismember grep -n nginx 'username=liamg&password=.*'
```

### Search for a pattern across all processes
```bash
# find a github api token across all processes
dismember grep 'gh[pousr]_[0-9a-zA-Z]{36}'
```

### Render a complete process tree
```bash
# defaults to pid=1 to show all (accessible) processes
dismember tree
```
