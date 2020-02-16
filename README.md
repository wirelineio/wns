# Wireline Naming Service

Wireline Naming Service (WNS) is a custom blockchain built using Cosmos SDK.

## Getting Started

### Installation

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

### Initializing the Local Node

# TODO(burdon): Create single root directory (e.g., ~/.wns)
Initialize the blockchain if it has never been run before (or run `rm -rf ~/.wnsd ~/.wnscli` first to delete all existing data and start over).

Initialize configuration files and genesis file.

# TODO(burdon): Explain why the developer needs to pick a name.
# TODO(burdon): This seems cosmos specific -- a script should handle set-up.
# TODO(burdon): Dumps a scary JSON object (no idea what just happened).
```bash
$ wnsd init <NAME> --chain-id wireline
```

Change the staking token name in `~/.wnsd/config/genesis.json` from `stake` to `uwire`.

# TODO(burdon): Why does the developer need to do this?
```
    "staking": {
      "params": {
        "unbonding_time": "1814400000000000",
        "max_validators": 100,
        "max_entries": 7,
        "bond_denom": "stake"    # --------> Change from "stake" TO "uwire".
      }
    }
```

# TODO(burdon): Why "uwire"?
```bash
$ sed -i '' 's/stake/uwire/g' ~/.wnsd/config/genesis.json
```

Optionally, change the following parameters for local testing purposes to the desired value:

* `app_state.nameservice.params.record_rent` - Record rent per period.
* `app_state.nameservice.params.record_expiry_time` - Record expiry time in nanoseconds.
* `app_state.bond.params.max_bond_amount` - Maximum amount a bond can hold.
* `app_state.gov.voting_params.voting_period` - Voting period for governance proposals (e.g. param changes).

Create a genesis validator account provisioned with 100 million WIRE.

# TODO(burdon): Mention passpharse (here and below).
# TODO(burdon): Generate the mnemonic?
Use the following mnemonic (or pass your own saved mnemonic from earlier runs):
`salad portion potato insect unknown exile lion soft layer evolve flavor hollow emerge celery ankle sponsor easy effort flush furnace life maximum rotate apple`

NOTE: To generate a new mnemonic & key, skip the --recover option.

# TODO(burdon): Don't put comments in code blocks -- put them in the document.
```bash
$ wnscli keys add root --recover
$ wnsd add-genesis-account $(wnscli keys show root -a) 100000000000000uwire
```

# TODO(burdon): Explain what this means.
Optionally, create a `faucet` genesis account (note the mnemonic).

```bash
$ wnscli keys add faucet
$ wnsd add-genesis-account $(wnscli keys show faucet -a) 100000000000000uwire
```

# TODO(burdon): I'm being asked for multiple passphrases for things I don't understand and multiple mnemonics.
# Guarantee this will cost you hours supporting the team in trying to reset a week after they've gone through this.

# TODO(burdon): Why is this a user-step? Could I call the chain something else?
Configure the CLI:


# TODO(burdon): /Users/burdon/.wnscli/config/config.toml does not exist (sounds like an error).
# TODO(burdon): If we have control over this CLI then silent if OK (with --verbose option)

```bash
$ wnscli config chain-id wireline
$ wnscli config output json
$ wnscli config indent true
$ wnscli config trust-node true
```

# TODO(burdon): Explain (not part of config).
# Validator stake/bond => 10 million WIRE (out of total 100 million WIRE).
# TODO(burdon): Here I'm asked for a "password" where previously I've been asked for a "passphrase"
# TODO(burdon): Commands below dump very large JSON objects.
```bash
$ wnsd gentx --name root --amount 10000000000000uwire

$ wnsd collect-gentxs
$ wnsd validate-genesis
```

### Starting the Node

# TODO(burdon): Common terminology (node?)
Start the server.

```bash
$ wnsd start start --gql-server --gql-playground
```

# TODO(burdon): Threw an exception:
```I[2020-02-16|12:43:53.640] Starting ABCI with Tendermint                module=main
panic: [{"msg_index":0,"success":false,"log":"{\"codespace\":\"staking\",\"code\":102,\"message\":\"invalid coin denomination\"}"}]```


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
