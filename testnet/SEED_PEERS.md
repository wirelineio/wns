# Seed Peers

Update `~/.wnsd/config/config.toml` with:

```text
seeds="06c238cf0366a3a3b83d1f0905fa186d4abe6e32@192.168.0.15:26656"
```

## Adding Seed Nodes

Get the Tendermint Node ID of the new peer:

```bash
$ wnsd tendermint show-node-id
```

Get the public hostname/IP for the machine to add as a new peer and update (`<node-id>@<host/IP>:26656`) the above list. Peers are separated by commas.


