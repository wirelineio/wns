#!/bin/sh

# TODO(ashwin): To switch to https, we have to distribute self-signed cert to clients.
export WIRE_WNS_ENDPOINT="http://node1.dxos.network:9473/graphql"

if [[ -z "$WIRE_WNS_USER_KEY" ]]; then
  echo "WNS user key not found. Set WIRE_WNS_USER_KEY and try again."
  return
fi

export WIRE_WNS_BOND_ID=
NUM_BONDS=$(wire wns bond list | jq -e ". | length")
if [ "$NUM_BONDS" -eq "0" ]; then
  wire wns bond create --type uwire --quantity 10000000000 > /dev/null
fi

export WIRE_WNS_BOND_ID=$(wire wns bond list | jq -r ".[0].id")

echo WIRE_WNS_ENDPOINT=${WIRE_WNS_ENDPOINT}
echo WIRE_WNS_USER_KEY=${WIRE_WNS_USER_KEY}
echo WIRE_WNS_BOND_ID=${WIRE_WNS_BOND_ID}
