# Wireline Naming Service

Wireline Naming Service (WNS) is a custom blockchain built using Cosmos SDK.

## Getting Started

### Setup Machine

* [Install golang](https://golang.org/doc/install) 1.13.0+ for the required platform.
* Test that `golang` has been successfully installed on the machine.

```bash
$ go version
go version go1.13 linux/amd64
```

Adding some ENV variables is necessary if `go mod` has never been used on the machine.

```bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.profile
echo "export GOBIN=\$GOPATH/bin" >> ~/.profile
echo "export PATH=\$PATH:\$GOBIN" >> ~/.profile
echo "export GO111MODULE=on" >> ~/.profile
source ~/.profile
```

Clone the repo (e.g. inside ~/wireline), build and install the binaries.

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

### Initialize Blockchain

Initialize the blockchain if it has never been run before (or run `rm -rf ~/.wnsd ~/.wnscli` first to delete all existing data and start over).

Initialize configuration files and genesis file.

```bash
# `my-node` is the name of the node.
$ wnsd init my-node --chain-id wireline
```

Change the staking token name in `~/.wnsd/config/genesis.json` from `stake` to `wire`.

```
    "staking": {
      "params": {
        "unbonding_time": "1814400000000000",
        "max_validators": 100,
        "max_entries": 7,
        "bond_denom": "stake"    # --------> Change from "stake" TO "wire".
      }
    }
```

```bash
$ sed -i '' 's/stake/wire/g' ~/.wnsd/config/genesis.json
```

```bash
$ wnscli keys add root --recover
# Use the following mnemonic (or pass your own saved mnemonic from earlier runs):
# salad portion potato insect unknown exile lion soft layer evolve flavor hollow emerge celery ankle sponsor easy effort flush furnace life maximum rotate apple
# To generate a new mnemonic & key, skip the --recover option.

# Create a genesis validator account provisioned with 100 million WIRE.
$ wnsd add-genesis-account $(wnscli keys show root -a) 100000000wire

# Optionally, create a `faucet` genesis account (note the mnemonic).
$ wnscli keys add faucet
$ wnsd add-genesis-account $(wnscli keys show faucet -a) 100000000wire
```

Configure the CLI to eliminate the need for the `chain-id` flag.

```bash
$ wnscli config chain-id wireline
$ wnscli config output json
$ wnscli config indent true
$ wnscli config trust-node true

# Validator stake/bond => 10 million WIRE (out of total 100 million WIRE).
$ wnsd gentx --name root --amount 10000000wire

$ wnsd collect-gentxs
$ wnsd validate-genesis
```

### Start Blockchain

Start the server.

```bash
$ wnsd start start --gql-server --gql-playground
```

Check that WNS is up and running by querying the GQL endpoint in another terminal.

```bash
$ curl -s -X POST -H "Content-Type: application/json" \
  -d '{ "query": "{ getStatus { version } }" }' http://localhost:9473/graphql | jq
```

## GQL Server API

The GQL server is controlled using the following `wnsd` flags:

* `--gql-server` - Enable GQL server (Available at http://localhost:9473/graphql).
* `--gql-playground` - Enable GQL playground app (Available at http://localhost:9473/console).
* `--gql-port` - Port to run the GQL server on (default 9473).

See `wnsd/x/nameservice/gql/schema.graphql` for the GQL schema.

## WNS CLI

[WNS CLI](https://github.com/wirelineio/registry-cli) provides commands within the `wire` utility for publishing and querying WNS records. 

## Testnets

### Development

Endpoints

* GQL: https://wns-testnet.dev.wireline.ninja/graphql
* GQL Playground: https://wns-testnet.dev.wireline.ninja/console
* RPC: tcp://wns-testnet.dev.wireline.ninja:26657

### Production

Endpoints

* GQL: https://wns-testnet.wireline.ninja/graphql
* GQL Playground: https://wns-testnet.wireline.ninja/console
* RPC: tcp://wns-testnet.wireline.ninja:26657

Note: The `wnscli` command accepts a `--node` flag for the RPC endpoint.

## Faucet

See https://github.com/wirelineio/faucet#environments.

## References

* https://golang.org/doc/install
* https://github.com/cosmos/cosmos-sdk
* https://cosmos.network/docs/tutorial/
* https://github.com/cosmos/sdk-application-tutorial
