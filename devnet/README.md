# Devnet

## Requirements

* [Hardware](https://github.com/dxos/xbox/blob/master/docs/hardware.md)
* [Ubuntu server](./SERVER.md)
* [Network](./NETWORK.md)

Note: These are instructions for setting up a validator node from genesis. To run a full node connected to the already running devnet, see the full node [setup](./FULL-NODE.md) doc.

## Endpoints

* GQL API: http://wns1.bozemanpass.net:9473/graphql
* GQL Console: http://wns1.bozemanpass.net:9473/console

## Validator Account Setup

Note: Run this step on every validator node.

Set an ENV variable with the mnemonic to be used for generating the validator account keys. Use an existing one generated earlier or create a new one using `wire keys generate`.

The mnemonic will be saved to `~/.wire/secrets` by the setup process, but also copy it to another safe location. There is no way to recover the account and associated funds if this mnemonic is lost.

```bash
$ export MNEMONIC="<MNEMONIC>"
```

Run the setup script (reset the node if prompted).

```bash
$ cd wns
$ ./scripts/setup.sh
```

Check-in the genesis transaction file created in `~/.wire/wnsd/config/gentx` to the `wns/devnet/gentx` folder.

Get the validator account address.

```bash
$ wnscli keys show root -a
cosmos174hrcf4x9nhwzt82qwns65esa0a7u05425jftp
```

Update SEED_ACCOUNTS.md with a new entry (validator address as above):

Note: Do NOT run this command, only copy it to the above file.

```text
wnsd add-genesis-account cosmos174hrcf4x9nhwzt82qwns65esa0a7u05425jftp 100000000000000uwire
```

## Genesis File Generation

Note: Run this step only on the initial validator node, to generate the consolidated `genesis.json` file.

Run the above setup. Delete existing contents in `~/.wire/wnsd/config/gentx` folder and copy all the gentx files from the repo to `~/.wire/wnsd/config/gentx`.

```bash
$ rm ~/.wire/wnsd/config/gentx/*
$ cp devnet/gentx/* ~/.wire/wnsd/config/gentx
```

Add the genesis accounts from [SEED_ACCOUNTS.md](./SEED_ACCOUNTS.md).

Re-generate the genesis.json file.

```bash
$ wnsd collect-gentxs
$ wnsd validate-genesis
```

Check-in the updated `~/.wire/wnsd/config/genesis.json` file to `wns/devnet/genesis.json`.

```bash
$ cp ~/.wire/wnsd/config/genesis.json devnet/genesis.json
```

## Genesis File Update

Note: Run this step on every validator node.

All validators should replace their `~/.wire/wnsd/config/genesis.json` file with the one in the repo, after the consolidated genesis.json has been generated.

```bash
$ cp devnet/genesis.json ~/.wire/wnsd/config/genesis.json
```

## Peer Setup

Note: Run this step on every validator node.

See [PEERS.md](./PEERS.md) to configure your node with peers. Once peers have been setup, the node can be started.

```bash
$ ./scripts/server.sh start
```

The devnet will generate blocks once 2/3 of voting power is online.

Check WNS logs to verify blocks are being generated.

```bash
$ ./scripts/server.sh log
```

## Post Setup

Note: To be run on each node once the devnet is operational and blocks are being generated.

Update `~/.profile` with account private key and node address. The private key can be looked up in `~/.wire/secrets`.

```
export WIRE_WNS_ENDPOINT="http://localhost:9473/graphql"
export WIRE_WNS_USER_KEY="<PRIVATE KEY>"
```

Apply the changes to `~/.profile`.

```bash
$ source ~/.profile
```

Generate a bond ID, which is required to pay for records.

```bash
$ wire wns create-bond --type uwire --quantity 10000000000
{
    "submit": "9F3E05DECE29D1B20F8148B8AEDA31058094036C2971ACA963A6ABE83A59587E"
}
```

List the bond IDs by owner address (`~/.wire/secrets` has the address to use).

```bash
$ wire wns list-bonds --owner cosmos1np8f3zzu6xss0m2rh2k7ugawegw0x29gh9n2lq
[
    {
        "id": "36550248ce6bd1d391825bc9111956dd899ac3ca03c238a20f79be49c8a9f806",
        "owner": "cosmos1np8f3zzu6xss0m2rh2k7ugawegw0x29gh9n2lq",
        "balance": [
            {
                "type": "uwire",
                "quantity": "10000000000"
            }
        ]
    }
]
```

Update `~/.profile` with the bond ID.

```
export WIRE_WNS_BOND_ID="<BOND ID>"
```

Apply the changes to `~/.profile`.

```bash
$ source ~/.profile
```


## Endpoints

### WNS

* GQL API: http://wns1.bozemanpass.net:9473/graphql
* GQL Playground: http://wns1.bozemanpass.net:9473/console
* TODO(ashwin): RPC Endpoint

### WNS Lite

* TODO(ashwin): Lite GQL API
* TODO(ashwin): Lite GQL Playground

### Faucet

* TODO(ashwin): API
* TODO(ashwin): GQL Playground

