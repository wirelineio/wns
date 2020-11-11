---
title: Devnet Moon
description: Instructions for setting up a devnet-moon full-node and optionally upgrading it to a validator node.
---

## Setup

Build and install the binaries:

```bash
$ cd wns
$ git checkout release-moon
$ make install
```

## Full Node Setup

```bash
$ export VALIDATOR_NAME="<NAME>"
$ ./scripts/setup.sh --reset --chain-id devnet-2 --node-name $VALIDATOR_NAME
$ cp networks/devnet-moon/genesis.json ~/.wire/wnsd/config/genesis.json
$ wnscli keys add $VALIDATOR_NAME
```

Request funds into the above validator account.

Update `~/.wire/wnsd/config/config.toml` with:

```text
persistent_peers = "<node-id>@<ip-address>:26656"
```

Optionally, update the above to include the new validator peer (Get node ID using `wnsd tendermint show-node-id`).

Start the node:

```bash
$ ./scripts/server.sh start --tail
```

## Full Node to Validator Node Upgrade

```bash
$ export VALIDATOR_NAME="<NAME>"

$ wnscli tx staking create-validator \
    --moniker $VALIDATOR_NAME \
    --chain-id "devnet-2" \
    --amount 10000000000000uwire \
    --pubkey $(wnsd tendermint show-validator) \
    --commission-max-change-rate "0.01" \
    --commission-max-rate "0.20" \
    --commission-rate "0.10" \
    --min-self-delegation "1" \
    --from $VALIDATOR_NAME
```

Check that the validator address is present in the latest validator set:

```bash
$ wnsd tendermint show-validator
$ wnscli query tendermint-validator-set
```

## Genesis File Creation

Note: To be run on a single machine.

```bash
$ ./scripts/setup.sh --reset --chain-id devnet-2
$ cp networks/devnet-moon/export.json ~/.wire/wnsd/config/genesis.json
$ wnsd collect-gentxs
$ wnsd validate-genesis
$ cp ~/.wire/wnsd/config/genesis.json networks/devnet-moon/genesis.json
```

Commit and push the file to GitHub.