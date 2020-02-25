#!/bin/bash

LOG=/tmp/wns.log
API_ENDPOINT=http://localhost:9473/graphql

function start_server ()
{
  stop_server
  set -x

  rm -f ${LOG}

  # Start the server.
  nohup wnsd start --gql-server --gql-playground > ${LOG} 2>&1 &

  if [[ $1 = "--tail" ]]; then
    log
  fi
}

function stop_server ()
{
  set -x
  killall wnsd
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

select yn in "start" "stop" "log" "test"; do
  command $yn
  exit
done
