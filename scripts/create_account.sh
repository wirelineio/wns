#!/bin/sh

TMP_ACCOUNT_FILE="/tmp/account.json"

ENDPOINT=
USER_KEY=
VERBOSE=

# Parse options and positional args.
POSITIONAL=()
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    -v|--verbose)
    VERBOSE=1
    shift
    ;;
    --endpoint)
    ENDPOINT="$2"
    shift
    shift
    ;;
    --user-key)
    USER_KEY="$2"
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

if [[ -z "${ENDPOINT}" ]] || [[ -z "${USER_KEY}" ]]; then
  echo "Usage: ./scripts/create_account.sh --endpoint <ENDPOINT> --user-key <PRIVATE KEY>"
  exit
fi

wire keys generate --json > "${TMP_ACCOUNT_FILE}"

NEW_ACCOUNT_ADDRESS=`jq -r ".address" ${TMP_ACCOUNT_FILE}`
NEW_ACCOUNT_USER_KEY=`jq -r ".privateKey" ${TMP_ACCOUNT_FILE}`

# Send tokens to account.
wire wns tokens send --address "${NEW_ACCOUNT_ADDRESS}" --type uwire --quantity 100000000000 --user-key "${USER_KEY}" --endpoint "${ENDPOINT}" > /dev/null

# Create bond for account.
wire wns bond create --type uwire --quantity 10000000000 --user-key "${NEW_ACCOUNT_USER_KEY}" --endpoint "${ENDPOINT}" > /dev/null
NEW_ACCOUNT_BOND_ID=$(wire wns bond list --owner ${NEW_ACCOUNT_ADDRESS} --endpoint "${ENDPOINT}" | jq -r ".[0].id")

if [[ ! -z "${VERBOSE}" ]]; then
  wire wns account get --address "${NEW_ACCOUNT_ADDRESS}" --endpoint "${ENDPOINT}"
  wire wns bond list --owner "${NEW_ACCOUNT_ADDRESS}" --endpoint "${ENDPOINT}"
fi

echo "-------------------------------------------------------------------------------------------------------------"
cat "${TMP_ACCOUNT_FILE}"
echo ""
echo "export WIRE_WNS_ENDPOINT=${ENDPOINT}"
echo "export WIRE_WNS_USER_KEY=${NEW_ACCOUNT_USER_KEY}"
echo "export WIRE_WNS_BOND_ID=${NEW_ACCOUNT_BOND_ID}"
echo "-------------------------------------------------------------------------------------------------------------"
