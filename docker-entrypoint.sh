#!/bin/bash
set -euo pipefail
readonly IFS=$'\n\t'

# If there are any arguments then we want to run those instead
if [[ "$1" == "-"* || -z $1 ]]; then
    exec /opt/aws-coreos-dashboard "$@"
else
    exec "$@"
fi
