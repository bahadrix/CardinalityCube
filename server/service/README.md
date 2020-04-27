# Command Line Interface Tool for Cardinality Cube Server

## Dependencies

Following packages needed by GOCMZQ

- pkg-config
- libczmq-dev
- libsodium-dev

You can install them in Debian based system like:
```bash
apt-get install -y pkg-config libczmq4 libczmq-dev libsodium-dev
```


## Tool Usage
```
Usage:
  server [command]

Available Commands:
  help        Help about any command
  start       Start a Cardinality Cube server

Flags:
      --config string   config file (default is $HOME/.ccserver.yaml)
  -h, --help            help for server
  -t, --toggle          Help message for toggle

Use "server [command] --help" for more information about a command.

```
