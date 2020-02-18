#!/bin/sh

#
# Initial set-up.
#

WNS_CLI_CONFIG_DIR=${HOME}/.wnscli
WNS_SERVER_CONFIG_DIR=${HOME}/.wnsd

WNS_CLI_EXTRA_ARGS="--home ${WNS_CLI_CONFIG_DIR}"
WNS_SERVER_EXTRA_ARGS="--home ${WNS_SERVER_CONFIG_DIR}"

NODE_NAME=WIRELINE
CHAIN_ID=wireline
DENOM=uwire

# TODO(burdon): Generate and save in ~/.wireline/secrets?
MNEMONIC="salad portion potato insect unknown exile lion soft layer evolve flavor hollow emerge celery ankle sponsor easy effort flush furnace life maximum rotate apple"

# TODO(ashwin): Save to ~/.wireline/secrets?
PASSPHRASE="temp12345"

function reset ()
{
  killall -SIGKILL wnsd
  rm -rf ${WNS_SERVER_CONFIG_DIR}
  rm -rf ${WNS_CLI_CONFIG_DIR}
}

function init_config ()
{
  wnscli config chain-id ${CHAIN_ID} $WNS_CLI_EXTRA_ARGS
  wnscli config output json $WNS_CLI_EXTRA_ARGS
  wnscli config indent true $WNS_CLI_EXTRA_ARGS
  wnscli config trust-node true $WNS_CLI_EXTRA_ARGS
}

function init_node ()
{
  # Init the chain.
  wnsd init ${NODE_NAME} --chain-id ${CHAIN_ID} $WNS_SERVER_EXTRA_ARGS

  # Change the staking unit.
  sed -i '' "s/stake/${DENOM}/g" "${WNS_SERVER_CONFIG_DIR}/config/genesis.json"

  # Change max bond amount from 10wire to 1000wire for easier local testing.
  sed -i '' 's/10wire/10000wire/g' "${WNS_SERVER_CONFIG_DIR}/config/genesis.json"
}

function init_root ()
{
  # Create the root account.
  echo "${PASSPHRASE}\n${MNEMONIC}" | wnscli keys add root --recover $WNS_CLI_EXTRA_ARGS
  wnsd add-genesis-account $(wnscli keys show root -a $WNS_CLI_EXTRA_ARGS) 100000000000000uwire $WNS_SERVER_EXTRA_ARGS
  echo "$PASSPHRASE" | wnsd gentx --name root --amount 10000000000000uwire $WNS_SERVER_EXTRA_ARGS --home-client ${WNS_CLI_CONFIG_DIR}
  wnsd collect-gentxs $WNS_SERVER_EXTRA_ARGS
  wnsd validate-genesis $WNS_SERVER_EXTRA_ARGS
}

#
# Options
#

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

init_config
init_node
init_root

echo "OK"
