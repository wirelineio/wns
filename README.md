# Wireline Naming Service

The Wireline Naming Service (WNS) is a custom blockchain built using Cosmos SDK.

## Getting Started

### Installation

* Install latest `wire` CLI before setting up `wnsd`.
* [Install golang](https://golang.org/doc/install) 1.13.0+ for the required platform.
* Test that `golang` has been successfully installed on the machine.

```bash
$ go version
go version go1.13 linux/amd64
```

Set the followin ENV variables (if `go mod` has never been used on the machine).

```bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.profile
echo "export GOBIN=\$GOPATH/bin" >> ~/.profile
echo "export PATH=\$PATH:\$GOBIN" >> ~/.profile
echo "export GO111MODULE=on" >> ~/.profile
source ~/.profile
```

Clone the repo then build and install the binaries.

```bash
$ cd ~/wireline
$ git clone git@github.com:wirelineio/wns.git
$ cd wns
$ make install
```

Test that the following commands work:

```bash
$ wnsd help
$ wnscli help
```

### Initializing the Local Node

```bash
$ ./scripts/setup.sh
```

### Working with the Local Node

Start the node:

```bash
$ ./scripts/server.sh start
```

Test if the node is up:

```bash
$ ./scripts/server.sh test
```

View the logs:

```bash
$ ./scripts/server.sh log
```

Stop the node:

```bash
$ ./scripts/server.sh stop
```


## WNS CLI

`wire` CLI provides [commands](https://github.com/wirelineio/incubator/blob/master/dxos/wns-cli/README.md) for publishing and querying WNS records.

## Tests

See https://github.com/wirelineio/registry-client#tests


## GQL Server API

The GQL server is controlled using the following `wnsd` flags:

* `--gql-server` - Enable GQL server (Available at http://localhost:9473/graphql).
* `--gql-playground` - Enable GQL playground app (Available at http://localhost:9473/console).
* `--gql-port` - Port to run the GQL server on (default 9473).

See `wnsd/gql/schema.graphql` for the GQL schema.


## References

* https://golang.org/doc/install
* https://github.com/cosmos/cosmos-sdk
* https://cosmos.network/docs/tutorial/
* https://github.com/cosmos/sdk-application-tutorial
