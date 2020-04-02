# Peers

Update `~/.wire/wnsd/config/config.toml` with:

```text
persistent_peers = "cf258070947018977f8d6f8a16b6f9f0c1cd0fdb@139.178.68.130:26656,9380b3f500ae36d256a8e3f6b5a81ceb62768308@139.178.68.131:26656,57706d1ce3f23d22da8eb958c7cd148544c1b8ae@wns1.deepstacksoft.com:26656"
```

## Adding Peers

Get the Tendermint Node ID of the new peer:

```bash
$ wnsd tendermint show-node-id
```

Get the public hostname/IP for the machine to add as a new peer and update (`<node-id>@<host/IP>:26656`) the above list. Peers are separated by commas.


## Troubleshooting

* If the node hostname/IP is not routable, nodes might have trouble connecting to each other. Try setting `addr_book_strict = false` in `~/.wire/wnsd/config/config.toml`.
* If the node does not have a static IP, [reverse port forwarding](./NETWORK.md) can be used to tunnel through a remote machine that has a public IP/hostname.
