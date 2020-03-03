# Peers

Update `~/.wireline/wnsd/config/config.toml` with:

```text
persistent_peers="b39acbd2773e0217a308abe4bab1b065fc58f55f@192.168.0.15:26656,b6e02ace3f42316ad62b67f19692c3941e37824e@192.168.0.42:26656,eba419c642c53eb37378c5045ecaf04e316b505d@192.168.0.86:26656"
```

## Adding Peers

Get the Tendermint Node ID of the new peer:

```bash
$ wnsd tendermint show-node-id
```

Get the public hostname/IP for the machine to add as a new peer and update (`<node-id>@<host/IP>:26656`) the above list. Peers are separated by commas.


## Troubleshooting

* If node hostname/IP is not routable, nodes might have trouble connecting to each other. Try setting `addr_book_strict = false` in `~/.wireline/wnsd/config/config.toml`.