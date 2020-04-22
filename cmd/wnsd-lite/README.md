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

To enable the lite node to periodically discover additional RPC endpoints from WNS, pass it a GQL API endpoint (in the example below, it's the lite node endpoint itself).

```bash
$ ./scripts/lite/server.sh start --node "tcp://node1.dxos.network:26657" --endpoint "http://127.0.0.1:9475/api"
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

### RPC Endpoint Discovery

Currently, RPC endpoints are discovered by querying for `xbox` type records with a `wns.rpc` field.

To register a `xbox` with a WNS RPC endpoint:

```bash
$ wire xbox register --id 'wrn:xbox:ashwinp/wns1' --version 0.0.1 --data.wns.rpc='tcp://45.79.120.249:26657'
```

```
$ wire wns record get --id Qmae4rq7QzLwz4qrqoDHa29w3CXQGSA8G766zLmM5yWVrU
[
  {
    "id": "Qmae4rq7QzLwz4qrqoDHa29w3CXQGSA8G766zLmM5yWVrU",
    "type": "wrn:xbox",
    "name": "ashwinp/wns1",
    "version": "0.0.1",
    "owners": [
      "233b436a205539f0f8082507e300fc5f3ca9eb0a"
    ],
    "bondId": "8a359128068c85f9982a36308772057d098f16dc21288e312205bdf60a6961e9",
    "createTime": "2020-04-22T07:58:31.839889941",
    "expiryTime": "2021-04-22T07:58:31.839889941",
    "attributes": {
      "type": "wrn:xbox",
      "version": "0.0.1",
      "wns": "{\"rpc\":\"tcp://45.79.120.249:26657\"}",
      "name": "ashwinp/wns1"
    }
  }
]
```
