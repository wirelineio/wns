#!/bin/bash

LOG=/tmp/wns-lite.log
GQL_SERVER_PORT=9475
API_ENDPOINT="http://localhost:${GQL_SERVER_PORT}/graphql"
WNS_NODE_ADDRESS=localhost

function start_server ()
{
  stop_server
  set -x

  rm -f ${LOG}

  if [[ ! -z "$1" ]]; then
    WNS_NODE_ADDRESS="$1"
  fi

  echo $WNS_NODE_ADDRESS

  # Start the server.
  nohup wnsd-lite start --gql-port ${GQL_SERVER_PORT} --node "tcp://${WNS_NODE_ADDRESS}:26657" > ${LOG} 2>&1 &
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

  tail -f ${LOG}
}

function test ()
{
  set -x
  curl -s -X POST -H "Content-Type: application/json" -d '{ "query": "{ getStatus { version } }" }' ${API_ENDPOINT} | jq
}

function command ()
{
  case $1 in
    start ) start_server $2; exit;;
    stop ) stop_server; exit;;
    log ) log; exit;;
    test ) test; exit;;
  esac
}

command=$1
if [[ ! -z "$command" ]]; then
  command $1 $2
  exit
fi

select oper in "start" "stop" "log" "test"; do
  command $oper
  exit
done
