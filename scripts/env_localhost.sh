#!/bin/sh

PASSPHRASE="12345678"

CHECK_WNS_RUNNING=`netstat -na | grep 9473`
if [ "$?" -eq "1" ]; then
  echo "WNS is not running. Start the server and retry."
else
  export WIRE_WNS_ENDPOINT="http://localhost:9473/graphql"
  export WIRE_WNS_USER_KEY="b1e4e95dd3e3294f15869b56697b5e3bdcaa24d9d0af1be9ee57d5a59457843a"
  export WIRE_WNS_BOND_ID=

  NUM_BONDS=$(wire wns list-bonds | jq -e ". | length")
  if [ "$NUM_BONDS" -eq "0" ]; then
    wire wns create-bond --type uwire --quantity 10000000000 > /dev/null
  fi

  export WIRE_WNS_BOND_ID=$(wire wns list-bonds | jq -r ".[0].id")

  echo WIRE_WNS_ENDPOINT=${WIRE_WNS_ENDPOINT}
  echo WIRE_WNS_USER_KEY=${WIRE_WNS_USER_KEY}
  echo WIRE_WNS_BOND_ID=${WIRE_WNS_BOND_ID}
fi
