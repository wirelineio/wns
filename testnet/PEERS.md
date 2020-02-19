# Peers

Update `~/.wnsd/config/config.toml` with:

```text
persistent_peers="8b43cc4b801609e280e605fe7814bf90cc7101a2@192.168.0.15:26656"
```

## Adding Seed Nodes

Get the Tendermint Node ID of the new peer:

```bash
$ wnsd tendermint show-node-id
```

Get the public hostname/IP for the machine to add as a new peer and update (`<node-id>@<host/IP>:26656`) the above list. Peers are separated by commas.


## Troubleshooting

* If node hostname/IP is not routable, nodes might have trouble connecting to each other. Try setting `addr_book_strict = false` in `~/.wnsd/config/config.toml`.