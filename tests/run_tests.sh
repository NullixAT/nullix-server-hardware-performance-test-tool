#!/usr/bin/bash
SCRIPTDIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

rm -rf $SCRIPTDIR/../nshptt/
go run $SCRIPTDIR/../tool.go --create-config
cp $SCRIPTDIR/config.json $SCRIPTDIR/../nshptt/config.json
go run $SCRIPTDIR/../tool.go --create-test-files
go run $SCRIPTDIR/../tool.go --run