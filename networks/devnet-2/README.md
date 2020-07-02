# DEVNET-2

Build and install the binaries:

```bash
$ cd wns
$ git checkout feature-mechanisms
$ make install
```

Choose a name for your validator (short alphanumeric string, without spaces):

```bash
$ export VALIDATOR_NAME='<NAME>'
```

Generate the node and validator keys and save the generated mnemonic to a safe place:

```bash
$ ./scripts/setup.sh --reset --chain-id devnet-2 --node-name $VALIDATOR_NAME
$ cp networks/devnet-2/genesis.json ~/.wire/wnsd/config/genesis.json
$ wnscli keys add $VALIDATOR_NAME
```

Generate a staking transaction (sign with your validator key):

```bash
$ wnsd add-genesis-account $(wnscli keys show $VALIDATOR_NAME -a) 100000000000000uwire
$ wnsd gentx --name $VALIDATOR_NAME --amount 10000000000000uwire --output-document $VALIDATOR_NAME.json
```

Commit the generated file to `wns/networks/devnet-2/gentx` folder and push the changes to `feature-mechanisms` branch.
