---
title: Devnet Moon
description: Instructions for setting up a devnet-moon full-node and optionally upgrading it to a validator node.
---

## Setup

Build and install the binaries:

```bash
$ cd wns
$ git pull
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
persistent_peers = "f615c77be9de710864e48d74717bb6d343a3e50a@172.105.37.214:26656,213ce5cfaed99146c738cfca971a4f3a1dfe6d22@139.178.68.131:26656,20161eff6d0b1a1f0f26d86b95b7d948739e1f00@139.178.68.130:26656"
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
