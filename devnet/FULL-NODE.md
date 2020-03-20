# Devnet Full Node Setup

Note: These are instructions for setting up a full node connected to the already running devnet. To run a validator from genesis, see the validator [setup](./README.md) doc.

## Full Node Account Setup

Set an ENV variable with the mnemonic to be used for generating the full node account keys. Use an existing one generated earlier or create a new one using `wire keys generate`.

The mnemonic will be saved to `~/.wireline/secrets` by the setup process, but also copy it to another safe location. There is no way to recover the account and associated funds if this mnemonic is lost.

```bash
$ export MNEMONIC="<MNEMONIC>"
```

Run the setup script (reset the node if prompted).

```bash
$ cd wns
$ ./scripts/setup.sh
```

## Genesis File Update

Replace `~/.wireline/wnsd/config/genesis.json` file with the one in the repo.

```bash
$ cp devnet/genesis.json ~/.wireline/wnsd/config/genesis.json
```

## Peer Setup

See [PEERS.md](./PEERS.md) for the value of `persistent_peers`, and update it as described. Skip the other sections. Once peers have been setup, the node can be started.

```bash
$ ./scripts/server.sh start
```

Note: On first run, the node needs to catch up to the existing devnet block height, so the output will scroll rather quickly.
