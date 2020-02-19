# Testnet

## Validator Account Setup

Set an ENV variable with the mnemonic to be used for generating the validator account keys. Use an existing one generated earlier or create a new one using `wire keys generate`.

```bash
$ export MNEMONIC="<MNEMONIC>"
```

Run the setup script (reset the node if required).

```bash
$ ./scripts/setup.sh
```

Check-in the genesis transaction file created in `~/.wnsd/config/gentx` to the `wns/testnet/gentx` folder.


Get the root account address.

```bash
$ wnscli keys show root -a
cosmos1hfz2f3wefu7pwrafdnu9pt5s7y0h924j66hld4
```

Update SEED_ACCOUNTS.md with a new entry (validator address as above).

```text
wnsd add-genesis-account cosmos1hfz2f3wefu7pwrafdnu9pt5s7y0h924j66hld4 100000000000000uwire
```


## Genesis JSON Generation

Run the above setup.

Copy all the gentx files to `~/.wnsd/config/gentx`.

Add the genesis accounts from SEED_ACCOUNTS.md.

Re-generate the genesis.json file.

```bash
$ wnsd collect-gentxs
$ wnsd validate-genesis
```

Check-in the updated `~/.wnsd/config/genesis.json` file to `wns/testnet/genesis.json`.

All validators should replace their `~/.wnsd/config/genesis.json` file with the one in the repo.


## Configure Seed Peers

See SEED_PEERS.md to configure your node with seed peers.

