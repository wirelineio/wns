# WNS

## Clear Remote WNS

To clear a remote WNS, the following information is required:

* The RPC endpoint of the remote WNS (e.g. see https://github.com/wirelineio/wns#testnets).
* The mnemonic for an account that has funds on the WNS.

The following example will work for https://wns-testnet.dev.wireline.ninja/console.

Create an account on a different machine (e.g. laptop/desktop), using the mnemonic for the remote `root` account.

```
$ wnscli keys add root-testnet-dev --recover
# Enter a passphrase for the new account, repeat it when prompted.
# Use the following mnemonic for recovery:
# salad portion potato insect unknown exile lion soft layer evolve flavor hollow emerge celery ankle sponsor easy effort flush furnace life maximum rotate apple
```

Clear the remote WNS using the following command:

```
$ wnscli tx nameservice clear --from root-testnet-dev --node tcp://wns-testnet.dev.wireline.ninja:26657
# Enter passphrase when prompted.
```

Use the GQL playground (https://wns-testnet.dev.wireline.ninja/console) to query and confirm that all records are gone.
