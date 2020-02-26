# Devnet

## Validator Account Setup

Set an ENV variable with the mnemonic to be used for generating the validator account keys. Use an existing one generated earlier or create a new one using `wire keys generate`.

```bash
$ export MNEMONIC="<MNEMONIC>"
```

Run the setup script (reset the node if required).

```bash
$ ./scripts/setup.sh
```

Check-in the genesis transaction file created in `~/.wireline/wnsd/config/gentx` to the `wns/testnet/gentx` folder.

Get the validator account address.

```bash
$ wnscli keys show root -a
cosmos1hfz2f3wefu7pwrafdnu9pt5s7y0h924j66hld4
```

Update SEED_ACCOUNTS.md with a new entry (validator address as above):

```text
wnsd add-genesis-account cosmos1hfz2f3wefu7pwrafdnu9pt5s7y0h924j66hld4 100000000000000uwire
```


## Genesis JSON Generation

Run the above setup.

Delete existing contents in `~/.wireline/wnsd/config/gentx` folder and copy all the gentx files from the repo to `~/.wireline/wnsd/config/gentx`.

```bash
$ rm ~/.wireline/wnsd/config/gentx/*
$ cp testnet/gentx/* ~/.wireline/wnsd/config/gentx
```

Add the genesis accounts from SEED_ACCOUNTS.md.

Re-generate the genesis.json file.

```bash
$ wnsd collect-gentxs
$ wnsd validate-genesis
```

Check-in the updated `~/.wireline/wnsd/config/genesis.json` file to `wns/testnet/genesis.json`.

```bash
$ cp ~/.wireline/wnsd/config/genesis.json testnet/genesis.json
```

All validators should replace their `~/.wireline/wnsd/config/genesis.json` file with the one in the repo.

## Peer Setup

See PEERS.md to configure your node with peers. Once peers have been setup, the node can be started.

The testnet will generate blocks once 2/3 of voting power is online.
