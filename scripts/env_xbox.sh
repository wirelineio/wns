#!/bin/sh

SCRIPT_DIR="$(dirname "$0")"
source "${SCRIPT_DIR}/env_devnet.sh"

# Now override WNS endpoint to point to xbox.
export WIRE_WNS_ENDPOINT="http://xbox.local:9473/graphql"

# TODO(ashwin): Remove once CLI config switching is in place.
export WIRE_IPFS_SERVER='http://xbox.local:5001'
export WIRE_IPFS_GATEWAY='http://xbox.local:8888/ipfs'

echo WIRE_IPFS_SERVER=${WIRE_IPFS_SERVER}
echo WIRE_IPFS_GATEWAY=${WIRE_IPFS_GATEWAY}
