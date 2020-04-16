# WNS Lite

WNS Lite is a light weight read-only cache of the WNS graph database.

## Getting Started

### Installation

Install [WNS](../../README.md), then test that the following command works:

```bash
$ wnsd-lite help
```

### Initializing the Lite Node

```bash
$ ./scripts/lite/setup.sh --node "<WNS RPC ENDPOINT>"
```

Example:

```bash
$ ./scripts/lite/setup.sh --node "tcp://node1.dxos.network:26657"
```

## Working with the Lite Node

Start the node:

```bash
$ ./scripts/lite/server.sh start --node "<WNS RPC ENDPOINT>"
```

Example:

```bash
$ ./scripts/lite/server.sh start --node "tcp://node1.dxos.network:26657"
```

Test if the node is up:

```bash
$ ./scripts/lite/server.sh test
```

View the logs:

```bash
$ ./scripts/lite/server.sh log
```

Stop the node:

```bash
$ ./scripts/lite/server.sh stop
```
