# Devnet

## Requirements

* [Hardware](https://github.com/dxos/xbox/blob/master/docs/hardware.md)
* Static public IP or [remote port forwarding](https://www.ssh.com/ssh/tunneling/example#remote-forwarding)
  * Ports to forward: 26656 (e.g. `ssh -nNT -vvv -R 26656:localhost:26656 wns.example.org`)

## Validator Account Setup

Note: Run this step on every validator node.

Set an ENV variable with the mnemonic to be used for generating the validator account keys. Use an existing one generated earlier or create a new one using `wire keys generate`.

```bash
$ export MNEMONIC="<MNEMONIC>"
```

Run the setup script (reset the node if required).

```bash
$ ./scripts/setup.sh
```

Check-in the genesis transaction file created in `~/.wireline/wnsd/config/gentx` to the `wns/devnet/gentx` folder.

Get the validator account address.

```bash
$ wnscli keys show root -a
cosmos1hfz2f3wefu7pwrafdnu9pt5s7y0h924j66hld4
```

Update SEED_ACCOUNTS.md with a new entry (validator address as above):

```text
$ wnsd add-genesis-account cosmos1hfz2f3wefu7pwrafdnu9pt5s7y0h924j66hld4 100000000000000uwire
```


## Genesis JSON Generation

Note: Run this step on a single validator node, to generate the consolidated `genesis.json` file.

Run the above setup. Delete existing contents in `~/.wireline/wnsd/config/gentx` folder and copy all the gentx files from the repo to `~/.wireline/wnsd/config/gentx`.

```bash
$ rm ~/.wireline/wnsd/config/gentx/*
$ cp devnet/gentx/* ~/.wireline/wnsd/config/gentx
```

Add the genesis accounts from SEED_ACCOUNTS.md.

Re-generate the genesis.json file.

```bash
$ wnsd collect-gentxs
$ wnsd validate-genesis
```

Check-in the updated `~/.wireline/wnsd/config/genesis.json` file to `wns/devnet/genesis.json`.

```bash
$ cp ~/.wireline/wnsd/config/genesis.json devnet/genesis.json
```

All validators should replace their `~/.wireline/wnsd/config/genesis.json` file with the one in the repo.


## Peer Setup

See PEERS.md to configure your node with peers. Once peers have been setup, the node can be started.

The devnet will generate blocks once 2/3 of voting power is online.
