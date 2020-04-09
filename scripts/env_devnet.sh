#!/bin/sh

export WIRE_WNS_ENDPOINT="https://node1.dxos.network/wns/graphql"

if [[ -z "$WIRE_WNS_USER_KEY" ]]; then
  echo "WNS user key not found. Export WIRE_WNS_USER_KEY and try again."
  return
fi

if [[ -z "$WIRE_WNS_BOND_ID" ]]; then
  echo "WNS bond ID not found. Export WIRE_WNS_BOND_ID and try again."
  echo "Create bond: wire wns bond create --type uwire --quantity 1000000000"
  echo "View bonds : wire wns bond list --owner <ADDRESS>"
  return
fi

echo WIRE_WNS_ENDPOINT=${WIRE_WNS_ENDPOINT}
echo WIRE_WNS_USER_KEY=${WIRE_WNS_USER_KEY}
echo WIRE_WNS_BOND_ID=${WIRE_WNS_BOND_ID}
