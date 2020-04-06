#!/bin/sh

# TODO(ashwin): To switch to https, we have to distribute self-signed cert to clients.
export WIRE_WNS_ENDPOINT="http://xbox.local:9473/graphql"
export WIRE_WNS_USER_KEY="b1e4e95dd3e3294f15869b56697b5e3bdcaa24d9d0af1be9ee57d5a59457843a"
export WIRE_WNS_BOND_ID=

NUM_BONDS=$(wire wns bond list | jq -e ". | length")
if [ "$NUM_BONDS" -eq "0" ]; then
wire wns bond create --type uwire --quantity 10000000000 > /dev/null
fi

export WIRE_WNS_BOND_ID=$(wire wns bond list | jq -r ".[0].id")

echo WIRE_WNS_ENDPOINT=${WIRE_WNS_ENDPOINT}
echo WIRE_WNS_USER_KEY=${WIRE_WNS_USER_KEY}
echo WIRE_WNS_BOND_ID=${WIRE_WNS_BOND_ID}

# TODO(ashwin): Remove once CLI config switching is in place.
export WIRE_IPFS_SERVER='http://xbox.local:5001'
export WIRE_IPFS_GATEWAY='http://xbox.local:8888/ipfs'

echo WIRE_IPFS_SERVER=${WIRE_IPFS_SERVER}
echo WIRE_IPFS_GATEWAY=${WIRE_IPFS_GATEWAY}
