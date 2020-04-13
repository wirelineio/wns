#!/bin/sh

export WIRE_WNS_ENDPOINT="https://node1.dxos.network/wns/graphql"
WIRE_FAUCET_ENDPOINT="https://node1.dxos.network/faucet/graphql"
WIRE_WNS_NETWORK="devnet"
WIRE_HOME=~/.wire
SECRETS_FILE="${WIRE_HOME}/devnet.secrets.json"

# Clean variables since we're sourcing the script in an existing shell.
MNEMONIC=
VERBOSE=

# Parse options and positional args.
POSITIONAL=()
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    --mnemonic)
    MNEMONIC="$2"
    shift
    shift
    ;;
    -v|--verbose)
    VERBOSE=1
    shift
    ;;
    *)
    POSITIONAL+=("$1")
    shift
    ;;
  esac
done
set -- "${POSITIONAL[@]}"

# Check if secrets file exists, else generate a new one.
if [[ ! -f "${SECRETS_FILE}" ]]; then
  if [[ ! -z ${MNEMONIC} ]]; then
    # Restore account from provided mnemonic.
    [[ ! -z "${VERBOSE}" ]] && echo "Restoring account from provided mnemonic."
    wire keys generate --mnemonic "${MNEMONIC}" --json > "${SECRETS_FILE}"
  else
    # Create account with a random mnemonic.
    echo "Generating account with new mnemonic (${SECRETS_FILE})."
    wire keys generate --json > "${SECRETS_FILE}"
    cat "${SECRETS_FILE}" | jq
  fi
else
  [[ ! -z "${VERBOSE}" ]] && echo "Loading existing account (${SECRETS_FILE})."
fi

# Load private key and address from secrets file.
export WIRE_WNS_USER_KEY=`jq -r ".privateKey" ${SECRETS_FILE}`
ADDRESS=`jq -r ".address" ${SECRETS_FILE}`

# Check if the account has bonds.
NUM_BONDS=$(wire wns bond list --owner ${ADDRESS} | jq -e ". | length")
if [ "$NUM_BONDS" -eq "0" ]; then
  [[ ! -z "${VERBOSE}" ]] && echo "No bonds found."

  # Prompt for Tweet URL.
  echo "Post a Tweet with text 'Fund ${ADDRESS}' and paste the URL below."
  read TWEET_URL

  # Request funds from faucet.
  echo "Requesting funds from faucet."
  wire faucet request --faucet-endpoint ${WIRE_FAUCET_ENDPOINT} --post-url "${TWEET_URL}" > /dev/null

  sleep 10

  # Create a new bond.
  echo "Creating a new bond."
  wire wns bond create --type uwire --quantity 1000000000 > /dev/null
else
  [[ ! -z "${VERBOSE}" ]] && echo "Account has existing bonds."
fi

export WIRE_WNS_BOND_ID=$(wire wns bond list --owner ${ADDRESS} | jq -r ".[0].id")

echo WIRE_WNS_ENDPOINT=${WIRE_WNS_ENDPOINT}
echo WIRE_WNS_USER_KEY=${WIRE_WNS_USER_KEY}
echo WIRE_WNS_BOND_ID=${WIRE_WNS_BOND_ID}
