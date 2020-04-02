#!/bin/sh

# Note: These ENV vars are NOT used in the scripts, but are set to satisfy the CLI config checking.
# TODO(ashwin): Remove after CLI stops asking for config even for `wire keys generate`.
export WIRE_SIGNAL_ENDPOINT='http://localhost:4000'
export WIRE_ICE_ENDPOINTS='[]'
export WIRE_IPFS_SERVER='http://127.0.0.1:5001'
export WIRE_IPFS_GATEWAY='http://127.0.0.1:8001/ipfs'
# End of block.

export WIRE_WNS_ENDPOINT="http://localhost:9473/graphql"
export WIRE_WNS_USER_KEY="b1e4e95dd3e3294f15869b56697b5e3bdcaa24d9d0af1be9ee57d5a59457843a"
export WIRE_WNS_BOND_ID=

if [[ "$1" = "--skip-bond-id" ]]; then
  return
fi

CHECK_WNS_RUNNING=`netstat -na | grep 9473`
if [ "$?" -eq "1" ]; then
  echo "WNS is not running. Start the server and retry."
else
  NUM_BONDS=$(wire wns bond list | jq -e ". | length")
  if [ "$NUM_BONDS" -eq "0" ]; then
    wire wns bond create --type uwire --quantity 10000000000 > /dev/null
  fi

  export WIRE_WNS_BOND_ID=$(wire wns bond list | jq -r ".[0].id")

  echo WIRE_WNS_ENDPOINT=${WIRE_WNS_ENDPOINT}
  echo WIRE_WNS_USER_KEY=${WIRE_WNS_USER_KEY}
  echo WIRE_WNS_BOND_ID=${WIRE_WNS_BOND_ID}
fi