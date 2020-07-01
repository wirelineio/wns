# DEVNET-2

```bash
$ cd wns
$ git checkout feature-mechanisms
$ make install
$ ./scripts/setup.sh --reset --chain-id devnet-2
$ cp networks/devnet-2/genesis.json ~/.wire/wnsd/config/genesis.json
```

Generate the validator key and save the mnemonic to a safe place.

```bash
$ wnscli keys add validator
```

Generate a staking transaction (replace `<NAME>` with your name).

```bash
$ wnsd add-genesis-account $(wnscli keys show validator -a) 100000000000000uwire
$ wnsd gentx --name validator --amount 10000000000000uwire --output-document <NAME>.json
```

Commit the generated file to `wns/networks/devnet-2/gentx` folder and push the changes to `feature-mechanisms` branch.
