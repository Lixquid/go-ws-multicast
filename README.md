# ws-multicast

A simple websocket multicast server.

Any messages sent to the server are rebroadcast to every client *except the one that sent it*. Nothing more, nothing less.

## Configuration

ws-multicast supports the following flags:

- `-host` (Defaults to `localhost`)
    - The host to listen to. To allow applications on other devices / servers to access the multicast server, set this to `0.0.0.0`.
- `-post` (Defaults to `9994`)
    - The port to listen on.

## Development

Builds can be created with `make build`.

Tests can be run with `go test`.
