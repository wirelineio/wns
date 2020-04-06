# Devnet Account Setup

To publish records to the devnet, an account with sufficient funds is required.

## Creating an Account

```bash
$ wire keys generate
```

To use an existing mnemonic, pass it as a CLI option (`--mnemonic "<MNEMONIC>"`).

Copy the mnemonic to another safe location. There is no way to recover the account and associated funds if this mnemonic is lost.

## Funding the Account

Create a [Tweet](https://twitter.com/compose/tweet) with the account address in the text.

Request funds from the devnet faucet.

TODO(ashwin): Update devnet faucet endpoint once deployed.

```bash
$ wire faucet request --faucet-endpoint "<FAUCET ENDPOINT>" --post-url "<Tweet URL>"
```

Check that the account has received funds.

```bash
$ wire wns account get --address "<ADDRESS>" --endpoint http://wns1.bozemanpass.net:9473/graphql
```

Note: Request more funds by creating a new Tweet with the same address. The faucet has a configured limit per account.

## Creating a Bond

A bond is required to pay rent for records published to the devnet.

To create a bond automatically, source the devnet ENV script.

```bash
$ cd wns
$ source ./scripts/env_devnet.sh
```

Check bonds associated with the account (should not be empty).

```bash
$ wire wns bond list --owner "<ADDRESS>"
```