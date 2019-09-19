# Wireline Naming Service

```bash
rm -rf ~/.wnsd ~/.wnscli

# Initialize configuration files and genesis file
wnsd init my-node --chain-id wireline

wnscli keys add root --recover
# Use the following mnemonic:
# salad portion potato insect unknown exile lion soft layer evolve flavor hollow emerge celery ankle sponsor easy effort flush furnace life maximum rotate apple

wnsd add-genesis-account $(wnscli keys show root -a) 1000nametoken,100000000stake

# Configure your CLI to eliminate need for chain-id flag
wnscli config chain-id wireline
wnscli config output json
wnscli config indent true
wnscli config trust-node true

wnsd gentx --name root

wnsd collect-gentxs
wnsd validate-genesis

wnsd start
```