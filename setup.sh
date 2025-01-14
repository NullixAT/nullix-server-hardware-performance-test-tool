#!/usr/bin/bash
SCRIPTDIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

env GOBIN=$SCRIPTDIR/build go install github.com/tc-hib/go-winres@latest
chmod +x $SCRIPTDIR/build/go-winres