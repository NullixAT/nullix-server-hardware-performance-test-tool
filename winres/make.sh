#!/usr/bin/bash
SCRIPTDIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

$SCRIPTDIR/../build/go-winres/go-winres make