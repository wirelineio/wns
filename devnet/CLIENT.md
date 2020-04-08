# Devnet Client

Note: These are instructions to connect to an existing devnet trusted node as a client. To run a validator or full node, see the [setup](./README.md) doc instead.

## Endpoints

WNS

* GQL API: http://node1.dxos.network:9473/graphql , https://node1.dxos.network/wns/graphql
* GQL Console: http://node1.dxos.network:9473/console, https://node1.dxos.network/wns/console
* RPC Endpoint: tcp://node1.dxos.network:26657

WNS Lite

* GQL API: http://node1.dxos.network:9475/graphql , https://node1.dxos.network/wnslite/graphql
* GQL Console: http://node1.dxos.network:9475/console , https://node1.dxos.network/wnslite/console

Faucet

* GQL API: http://faucet.node1.dxos.network:4000/graphql , https://node1.dxos.network/faucet/graphql
* GQL Console: http://faucet.node1.dxos.network:4000/console , https://node1.dxos.network/faucet/console

## Working with the Devnet

### Querying

To query the devnet, update the config file (or use a command line flag) to connect to the above GQL API endpoint. No other changes are required.

### Publishing

To publish records, an [account](./ACCOUNT.md) needs to be setup.

Once the account is setup, the wire CLI can be used to registers records (e.g. app/bot).

## Configuration

TODO(ashwin): What's the recommended way to switch CLI between localhost, xbox.local and devnet?

To connect to the devnet, either

* Configure the CLI (`~/.wire/config`), or
* Export the private key for the devnet account (`export WIRE_WNS_USER_KEY="<PRIVATE KEY>"`), then run the following override script.

```bash
$ cd wns
$ source ./scripts/env_devnet.sh
```

## Troubleshooting

Ensure that the CLI is configured correctly or the following ENV variables are correct.

* WIRE_WNS_ENDPOINT - must be the above WNS GQL API endpoint
* WIRE_WNS_USER_KEY - must be the `privateKey` for the devnet account (from output of `wire keys generate`)
* WIRE_WNS_BOND_ID - must be a bond owned by the account, with sufficient funds (`wire wns bond list --owner <ACCOUNT ADDRESS>`)
