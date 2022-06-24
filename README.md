# Dismember

Dismember is a command-line toolkit for Linux that can be used to explore processes and (especially) their memory. Essentially for playing with `/proc`.

One core feature is the ability to scan the memory of all processes for common secrets, or for custom regular expressions.

![A gif showing dismember finding credentials from the memory of a browser](demo.gif)

Using the `grep` command, it can match a regular expression across all memory for all (accessible) processes. This could be used to find sensitive data in memory, identify a process by something included in its memory, or to interrogate a processes' memory for interesting information.

There are many built-in patterns included via the `scan` command, which effectively works as a secret scanner against the memory on your machine.

Dismember can be used to search memory of all processes it has access to, so running it as root is the most effective method.

Commands are also included to list processes, explore process status and related information, draw process trees, and more...

## Available Commands

| Command   | Description                                                                              | 
|-----------|------------------------------------------------------------------------------------------|
| `files`   | Show a list of files being accessed by a process                                         |
| `find`    | Find a PID given a process name. If multiple processes match, the first one is returned. |
| `grep`    | Search process memory for a given string or regex                                        |
| `info`    | Show information about a process                                                         |
| `kernel`  | Show information about the kernel                                                        | 
| `kill`    | Kill a process using SIGKILL                                                             | 
| `list`    | List all processes currently available on the system                                     | 
| `resume`  | Resume a suspended process using SIGCONT                                                 | 
| `scan`    | Search process memory for a set of predefined secret patterns                            | 
| `suspend` | Suspend a process using SIGSTOP (use 'dismember resume' to leave suspension)             | 
| `tree`    | Show a tree diagram of a process and all children (defaults to PID 1).                   | 

## Installation

Grab a binary from the [latest release](https://github.com/liamg/dismember/releases/latest) and add it to your path.

## Usage Examples

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

### Search for secrets in memory across all processes
```bash
# search all accessible memory for common secrets
dismember scan
```

## FAQ

> Isn't this information all just sitting in `/proc`?

Pretty much. Dismember just reads and presents it for the most part. If you can get away with `grep whatever /proc/[pid]/blah` then go for it! I built this as an educational experience because I couldn't sleep one night and stayed up late reading the `proc` man-pages (I live an extremely rock 'n' roll lifestyle). It's not a replacement for existing tools, but perhaps it can complement them.

> Do you know how horrific some of these commands seem when read out of context?

[Yes](https://twitter.com/liam_galvin/status/1540375769049960448).
