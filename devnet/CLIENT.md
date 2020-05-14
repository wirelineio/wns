# Devnet Client

Note: These are instructions to connect to an existing `devnet` trusted node as a client. To run a validator or full node, see the [setup](./README.md) doc instead.

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

To query the `devnet`, update the `wire` profile config file to the above GQL API endpoint. No other changes are required.

### Publishing

To publish records, an [account](./ACCOUNT.md) needs to be setup.

Once the account is setup, the `wire` CLI can be used to registers records (e.g. app/bot).

Activate the `devnet` CLI profile created during account setup.

## Troubleshooting

Ensure that the CLI profile is configured correctly.

```bash
$ export WIRE_PROFILE=<PROFILE NAME>
$ wire config
```

* `services.wns.server` - must be a valid `devnet` WNS endpoint
* `services.wns.userKey` - must be the `privateKey` for the `devnet` account
* `services.wns.bondId` - must be a bond owned by the account, with sufficient funds (`wire wns bond list --owner <ACCOUNT ADDRESS>`)
