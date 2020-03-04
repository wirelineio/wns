# Devnet

## Requirements

* [Hardware](https://github.com/dxos/xbox/blob/master/docs/hardware.md)
* Static public IP or [remote port forwarding](https://www.ssh.com/ssh/tunneling/example#remote-forwarding)
  * Ports to forward: 26656 (e.g. `ssh -nNT -vvv -R 26656:localhost:26656 wns.example.org`)
* [Ubuntu server setup](./SERVER.md)

## Validator Account Setup

Note: Run this step on every validator node.

Set an ENV variable with the mnemonic to be used for generating the validator account keys. Use an existing one generated earlier or create a new one using `wire keys generate`.

The mnemonic will be saved to `~/.wireline/secrets` by the setup process, but also save it to another safe location of your choice. There is no way to recover the account and associated funds if this mnemonic is lost.

```bash
$ export MNEMONIC="<MNEMONIC>"
```

Run the setup script (reset the node if required).

```bash
$ cd wns
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

Check WNS logs to verify blocks are being generated.

```bash
$ ./scripts/server.sh log
```

## Post Setup

Note: To be run on each node once the devnet is operational and blocks are being generated.

Update `~/.profile` with account private key and node address. The private key can be looked up in `~/.wireline/secrets`.

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

List the bond IDs by owner address (`~/.wireline/secrets` has the address to use).

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
export WIRE_WNS_BOND_ID="36550248ce6bd1d391825bc9111956dd899ac3ca03c238a20f79be49c8a9f806"
```

Apply the changes to `~/.profile`.

```bash
$ source ~/.profile
```
