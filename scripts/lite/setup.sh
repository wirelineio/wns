#!/bin/bash

#
# Initial set-up.
#

WNS_LITE_SERVER_CONFIG_DIR="${HOME}/.wireline/wnsd-lite"

CHAIN_ID=wireline

function reset ()
{
  killall -SIGTERM wnsd-lite
  rm -rf "${WNS_LITE_SERVER_CONFIG_DIR}"
}

function init_node ()
{
  if [[ ! -z "$1" ]]; then
    mkdir -p "${WNS_LITE_SERVER_CONFIG_DIR}/config"
    cp "$1" "${WNS_LITE_SERVER_CONFIG_DIR}/config/genesis.json"
  fi

  if [[ ! -z "$2" ]]; then
    EXTRA_ARGS="--height $2"
  fi

  # Init the node.
  wnsd-lite init --chain-id "${CHAIN_ID}" ${EXTRA_ARGS}
}

#
# Options
#

# Test if installed already.
if [[ -d "${WNS_LITE_SERVER_CONFIG_DIR}" ]]; then
  echo "Do you wish to RESET?"
  select yn in "Yes" "No"; do
    case $yn in
      Yes ) reset $1 $2; break;;
      No ) exit;;
    esac
  done
fi

init_node $1 $2

echo "OK"
