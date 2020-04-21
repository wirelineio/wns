#!/bin/bash

LOG="/tmp/wns-lite.log"
GQL_SERVER_PORT="9475"
GQL_PLAYGROUND_API_BASE=""
WNS_NODE_ADDRESS="tcp://localhost:26657"

POSITIONAL=()
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    --node)
    WNS_NODE_ADDRESS="$2"
    shift
    shift
    ;;
    --gql-port)
    GQL_SERVER_PORT="$2"
    shift
    shift
    ;;
    --gql-playground-api-base)
    GQL_PLAYGROUND_API_BASE="$2"
    shift
    shift
    ;;
    --log)
    LOG="$2"
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

function start_server ()
{
  stop_server
  set -x

  rm -f "${LOG}"

  # Start the server.
  nohup wnsd-lite start --gql-port "${GQL_SERVER_PORT}" --gql-playground-api-base "${GQL_PLAYGROUND_API_BASE}" --node "${WNS_NODE_ADDRESS}" --log-level debug > "${LOG}" 2>&1 &
}

function stop_server ()
{
  set -x
  killall wnsd-lite
}

function log ()
{
  echo
  echo "Log file: ${LOG}"
  echo

  tail -f "${LOG}"
}

function test ()
{
  set -x
  curl -s -X POST -H "Content-Type: application/json" -d '{ "query": "{ getStatus { version } }" }' "http://localhost:${GQL_SERVER_PORT}/graphql" | jq
}

function command ()
{
  case $1 in
    start ) start_server; exit;;
    stop ) stop_server; exit;;
    log ) log; exit;;
    test ) test; exit;;
  esac
}

command=$1
if [[ ! -z "$command" ]]; then
  command $1
  exit
fi

select oper in "start" "stop" "log" "test"; do
  command $oper
  exit
done
