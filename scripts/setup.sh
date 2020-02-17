#!/bin/sh

#
# Initial set-up.
#

WNS_CONFIG=${HOME}/.wnsd
WNS_CLI_CONFIG=${HOME}/.wnscli

NODE_NAME=WIRELINE
CHAIN_ID=wireline
DENOM=uwire

# TODO(burdon): Generate and save in ~/.wireline/secrets?
MNEMONIC="salad portion potato insect unknown exile lion soft layer evolve flavor hollow emerge celery ankle sponsor easy effort flush furnace life maximum rotate apple"

function reset ()
{
  rm -rf ${WNS_CONFIG}
  rm -rf ${WNS_CLI_CONFIG}
}

function init_config ()
{
  wnscli config chain-id ${CHAIN_ID}
  wnscli config output json
  wnscli config indent true
  wnscli config trust-node true
}

function init_node ()
{
  # Init the chain.
  wnsd init ${NODE_NAME} --chain-id ${CHAIN_ID}

  # Change the staking unit.
  sed -i '' "s/stake/${DENOM}/g" ~/.wnsd/config/genesis.json
}

function init_root ()
{
  echo
  echo "Use the dev mnemonic:"
  echo ${MNEMONIC}
  echo

  # Create the root account.
  wnscli keys add root --recover
  wnsd add-genesis-account $(wnscli keys show root -a) 100000000000000uwire
}

function init_account ()
{
  # Create a standard account via the faucet.
  wnscli keys add faucet
  wnsd add-genesis-account $(wnscli keys show faucet -a) 100000000000000uwire
  wnsd gentx --name root --amount 10000000000000uwire
  wnsd collect-gentxs
  wnsd validate-genesis
}

#
# Options
#

# Test if installed already.
if [[ -d "${WNS_CONFIG}" ]]; then
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
init_account

echo "OK"

