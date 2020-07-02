#!/bin/bash

#
# Initial set-up.
#

DEFAULT_MNEMONIC="salad portion potato insect unknown exile lion soft layer evolve flavor hollow emerge celery ankle sponsor easy effort flush furnace life maximum rotate apple"
DEFAULT_PASSPHRASE="12345678"

NODE_NAME=`hostname`
CHAIN_ID=wireline
DENOM=uwire

WNS_CLI_CONFIG_DIR="${HOME}/.wire/wnscli"
WNS_SERVER_CONFIG_DIR="${HOME}/.wire/wnsd"

WNS_CLI_EXTRA_ARGS="--home ${WNS_CLI_CONFIG_DIR}"
WNS_SERVER_EXTRA_ARGS="--home ${WNS_SERVER_CONFIG_DIR}"

POSITIONAL=()
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    --reset)
    RESET=1
    shift
    ;;
    --chain-id)
    CHAIN_ID="$2"
    shift
    shift
    ;;
    --node-name)
    NODE_NAME="$2"
    shift
    shift
    ;;
    --mnemonic)
    MNEMONIC="$2"
    shift
    shift
    ;;
    --passphrase)
    PASSPHRASE="$2"
    shift
    shift
    ;;
    *)
    POSITIONAL+=("$1")
    shift
    ;;
  esac
done
set -- "${POSITIONAL[@]}"

function init_secrets ()
{
  if [[ -z "${MNEMONIC}" ]]; then
    MNEMONIC="${DEFAULT_MNEMONIC}"
  fi

  if [[ -z "${PASSPHRASE}" ]]; then
    PASSPHRASE="${DEFAULT_PASSPHRASE}"
  fi
}

SED_ARGS=""

# On MacOS, sed needs `-i ''``. On Linux, just `-i`.
if [ "$(uname)" == "Darwin" ]; then
  SED_ARGS="''"
fi

function save_secrets ()
{
  mkdir -p ~/.wire
  echo "Root Account Mnemonic: ${MNEMONIC}" > ~/.wire/secrets
  echo "CLI Passphrase: ${PASSPHRASE}" >> ~/.wire/secrets
  echo "Wire CLI Keys:" >> ~/.wire/secrets
  wire keys generate --mnemonic="${MNEMONIC}" >> ~/.wire/secrets
}

function reset ()
{
  killall -SIGKILL wnsd
  rm -rf "${WNS_SERVER_CONFIG_DIR}"
  rm -rf "${WNS_CLI_CONFIG_DIR}"
}

function init_config ()
{
  # Configure the CLI to eliminate the need for the chain-id flag.
  wnscli config chain-id "${CHAIN_ID}"
  wnscli config output json
  wnscli config indent true
  wnscli config trust-node true
}

function init_node ()
{
  # Init the chain.
  wnsd init "${NODE_NAME}" --chain-id "${CHAIN_ID}"

  # Change the staking unit.
  sed -i $SED_ARGS "s/stake/${DENOM}/g" "${WNS_SERVER_CONFIG_DIR}/config/genesis.json"

  # Change max bond amount from 10wire to 1000wire for easier local testing.
  sed -i $SED_ARGS "s/10wire/10000wire/g" "${WNS_SERVER_CONFIG_DIR}/config/genesis.json"
}

function init_root ()
{
  # Create a genesis validator account provisioned with 100 million WIRE.
  echo -e "${PASSPHRASE}\n${MNEMONIC}" | wnscli keys add root --recover $WNS_CLI_EXTRA_ARGS
  wnsd add-genesis-account $(wnscli keys show root -a $WNS_CLI_EXTRA_ARGS) 100000000000000uwire $WNS_SERVER_EXTRA_ARGS

  # Validator stake/bond => 10 million WIRE (out of total 100 million WIRE).
  echo -e "${PASSPHRASE}" | wnsd gentx --name root --amount 10000000000000uwire $WNS_SERVER_EXTRA_ARGS --home-client $WNS_CLI_CONFIG_DIR
  wnsd collect-gentxs $WNS_SERVER_EXTRA_ARGS
  wnsd validate-genesis $WNS_SERVER_EXTRA_ARGS
}

#
# Options
#

if [[ ! -z "${RESET}" ]]; then
  reset
fi

# Test if installed already.
if [[ -d "${WNS_SERVER_CONFIG_DIR}" ]]; then
  echo "Do you wish to RESET?"
  select yn in "Yes" "No"; do
    case $yn in
      Yes ) reset; break;;
      No ) exit;;
    esac
  done
fi

init_secrets

init_config
init_node
init_root

save_secrets
