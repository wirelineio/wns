# Wireline Naming Service

Wireline Naming Service (WNS) is a custom blockchain built using Cosmos SDK.

## Getting Started

### Setup Machine

Install golang 1.13.0+ for your platform.

```
$ $ go version
go version go1.13 linux/amd64
```

Adding some ENV variables is necessary if you've never used `go mod` on your machine.

```
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bash_profile
echo "export GOBIN=\$GOPATH/bin" >> ~/.bash_profile
echo "export PATH=\$PATH:\$GOBIN" >> ~/.bash_profile
echo "export GO111MODULE=on" >> ~/.bash_profile
source ~/.bash_profile
```

Clone the repo (e.g. inside ~/wireline), build and install the binaries.

```
$ cd ~/wireline
$ git clone git@github.com:wirelineio/wns.git
$ cd wns
$ make install
```

Now you should be able to run the following commands:

```
$ wnsd help
$ wnscli help
```

### Initialize Blockchain

Initialize the blockchain if you're never run it before (or run `rm -rf ~/.wnsd ~/.wnscli` first to delete all existing data and start over).

Initialize configuration files and genesis file.

```bash
# `my-node` is the name of the node.
$ wnsd init my-node --chain-id wireline

$ wnscli keys add root --recover
# Use the following mnemonic:
# salad portion potato insect unknown exile lion soft layer evolve flavor hollow emerge celery ankle sponsor easy effort flush furnace life maximum rotate apple
# If you'd like to use your own mnemonic, don't pass the --recover option.

$ wnsd add-genesis-account $(wnscli keys show root -a) 100000000wire,100000000stake
```

Configure your CLI to eliminate need for chain-id flag.

```
$ wnscli config chain-id wireline
$ wnscli config output json
$ wnscli config indent true
$ wnscli config trust-node true

$ wnsd gentx --name root

$ wnsd collect-gentxs
$ wnsd validate-genesis

```

### Start Blockchain

Start the server.

```
$ wnsd start start --gql-server --gql-playground
```

Check that the WNS is up and running by querying the GQL endpoint in another terminal.

```
$ curl -s -X POST -H "Content-Type: application/json" \
  -d '{ "query": "{ getStatus { version } }" }' http://localhost:9473/query | jq
```

## GQL Server API

The GQL server is controlled using the following `wnsd` flags:

* `--gql-server` - Enable GQL server.
* `--gql-playground` - Enable GQL playground app (Available at http://localhost:9473/).
* `--gql-port` - Port to run the GQL server on (default 9473).

See `wnsd/x/nameservice/gql/schema.graphql` for the GQL schema.

## References

* https://golang.org/doc/install
* https://github.com/cosmos/cosmos-sdk
* https://cosmos.network/docs/tutorial/
* https://github.com/cosmos/sdk-application-tutorial
