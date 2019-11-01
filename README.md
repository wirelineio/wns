# Wireline Naming Service

Wireline Naming Service (WNS) is a custom blockchain built using Cosmos SDK.

## Getting Started

### Setup Machine

Install golang 1.13.0+ for your platform.

```
$ go version
go version go1.13 linux/amd64
```

Adding some ENV variables is necessary if you've never used `go mod` on your machine.

```
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.profile
echo "export GOBIN=\$GOPATH/bin" >> ~/.profile
echo "export PATH=\$PATH:\$GOBIN" >> ~/.profile
echo "export GO111MODULE=on" >> ~/.profile
source ~/.profile
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
  -d '{ "query": "{ getStatus { version } }" }' http://localhost:9473/graphql | jq
```

## GQL Server API

The GQL server is controlled using the following `wnsd` flags:

* `--gql-server` - Enable GQL server (Available at http://localhost:9473/graphql).
* `--gql-playground` - Enable GQL playground app (Available at http://localhost:9473/console).
* `--gql-port` - Port to run the GQL server on (default 9473).

See `wnsd/x/nameservice/gql/schema.graphql` for the GQL schema.

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

The testnets come with a genesis account (`root`) that can be used to transfer funds to a new account. Run these commands locally to restore the keys on your own machine.

Note: Access to the mnemonic means access to all funds in the account. Don't share or use this mnemonic for non-testing purposes.

```
$ wnscli keys add root-dev-env --recover

# Use the following mnemonic for recovery:
# salad portion potato insect unknown exile lion soft layer evolve flavor hollow emerge celery ankle sponsor easy effort flush furnace life maximum rotate apple

$ wnscli tx bank send $(wnscli keys show root-dev-env -a) cosmos1lpzffjhasv5qhn7rn6lks9u4dvpzpuj922tdmy 1000wire  --from root-dev-env --chain-id=wireline --node tcp://wns-testnet.dev.wireline.ninja:26657

# Replace cosmos1lpzffjhasv5qhn7rn6lks9u4dvpzpuj922tdmy with the address you want to transfer funds to.

# Query updated balances.

$ wnscli query account $(wnscli keys show root-dev-env -a) --chain-id=wireline --node tcp://wns-testnet.dev.wireline.ninja:26657
$ wnscli query account cosmos1lpzffjhasv5qhn7rn6lks9u4dvpzpuj922tdmy --chain-id=wireline --node tcp://wns-testnet.dev.wireline.ninja:26657
```

## References

* https://golang.org/doc/install
* https://github.com/cosmos/cosmos-sdk
* https://cosmos.network/docs/tutorial/
* https://github.com/cosmos/sdk-application-tutorial
