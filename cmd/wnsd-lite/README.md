# WNS Lite

WNS Lite is a light weight read-only cache of the WNS graph database.

## Getting Started

### Installation

Install [WNS](../../README.md), then test that the following command works:

```bash
$ wnsd-lite help
```

### Initializing the Lite Node

Export the current state from a WNS node and note the height at which the state was exported (available in the logs).

```bash
$ wnsd export > ./genesis.json
```

Initialize the lite node with the genesis.json file and the corresponding height.

```bash
$ ./scripts/lite/setup.sh ./genesis.json <height>
```

Example:

```bash
$ ./scripts/lite/setup.sh ./genesis.json 188
```

## Working with the Lite Node

Start the node:

```bash
$ ./scripts/lite/server.sh start
```

By default, the node connects to WNS running on localhost. To point it to another host, pass the hostname/IP as a command line argument.

```bash
$ ./scripts/lite/server.sh start <WNS hostname/IP>
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
