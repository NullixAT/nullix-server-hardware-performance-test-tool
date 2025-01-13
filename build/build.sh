#!/usr/bin/bash
SCRIPTDIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

rm $SCRIPTDIR/nshptt_*
rm $SCRIPTDIR/rsrc_*

$SCRIPTDIR/go-winres/go-winres make
cp $SCRIPTDIR/../rsrc_* $SCRIPTDIR

env GOOS=aix GOARCH=ppc64 go build -o $SCRIPTDIR/nshptt_aix_ppc64 $SCRIPTDIR/../tool.go 

# env GOOS=android GOARCH=amd64 go build -o $SCRIPTDIR/nshptt_android_amd64.exe $SCRIPTDIR/../tool.go  
# env GOOS=android GOARCH=386 go build -o $SCRIPTDIR/nshptt_android_386.exe $SCRIPTDIR/../tool.go  
# env GOOS=android GOARCH=arm go build -o $SCRIPTDIR/nshptt_android_arm.exe $SCRIPTDIR/../tool.go  
env GOOS=android GOARCH=arm64 go build -o $SCRIPTDIR/nshptt_android_arm64 $SCRIPTDIR/../tool.go  

env GOOS=darwin GOARCH=amd64 go build -o $SCRIPTDIR/nshptt_darwin_amd64 $SCRIPTDIR/../tool.go  
env GOOS=darwin GOARCH=arm64 go build -o $SCRIPTDIR/nshptt_darwin_arm64 $SCRIPTDIR/../tool.go  

# env GOOS=ios GOARCH=amd64 go build -o $SCRIPTDIR/nshptt_ios_amd64.exe $SCRIPTDIR/../tool.go  
# env GOOS=ios GOARCH=arm64 go build -o $SCRIPTDIR/nshptt_ios_arm64.exe $SCRIPTDIR/../tool.go  

env GOOS=windows GOARCH=amd64 go build -o $SCRIPTDIR/nshptt_win_amd64.exe $SCRIPTDIR/../tool.go && $SCRIPTDIR/go-winres/go-winres patch $SCRIPTDIR/nshptt_win_amd64.exe
env GOOS=windows GOARCH=386 go build -o $SCRIPTDIR/nshptt_win_386.exe $SCRIPTDIR/../tool.go  && $SCRIPTDIR/go-winres/go-winres patch $SCRIPTDIR/nshptt_win_386.exe
env GOOS=windows GOARCH=arm go build -o $SCRIPTDIR/nshptt_win_arm.exe $SCRIPTDIR/../tool.go  
env GOOS=windows GOARCH=arm64 go build -o $SCRIPTDIR/nshptt_win_arm64.exe $SCRIPTDIR/../tool.go  
rm $SCRIPTDIR/nshptt_win_*.bak

env GOOS=linux GOARCH=amd64 go build -o $SCRIPTDIR/nshptt_linux_amd64 $SCRIPTDIR/../tool.go  
env GOOS=linux GOARCH=386 go build -o $SCRIPTDIR/nshptt_linux_386 $SCRIPTDIR/../tool.go  
env GOOS=linux GOARCH=arm go build -o $SCRIPTDIR/nshptt_linux_arm $SCRIPTDIR/../tool.go  
env GOOS=linux GOARCH=arm64 go build -o $SCRIPTDIR/nshptt_linux_arm64 $SCRIPTDIR/../tool.go  

env GOOS=openbsd GOARCH=amd64 go build -o $SCRIPTDIR/nshptt_openbsd_amd64 $SCRIPTDIR/../tool.go  
env GOOS=openbsd GOARCH=386 go build -o $SCRIPTDIR/nshptt_openbsd_386 $SCRIPTDIR/../tool.go  
env GOOS=openbsd GOARCH=arm go build -o $SCRIPTDIR/nshptt_openbsd_arm $SCRIPTDIR/../tool.go  
env GOOS=openbsd GOARCH=arm64 go build -o $SCRIPTDIR/nshptt_openbsd_arm64 $SCRIPTDIR/../tool.go  

env GOOS=netbsd GOARCH=amd64 go build -o $SCRIPTDIR/nshptt_netbsd_amd64 $SCRIPTDIR/../tool.go  
env GOOS=netbsd GOARCH=386 go build -o $SCRIPTDIR/nshptt_netbsd_386 $SCRIPTDIR/../tool.go  
env GOOS=netbsd GOARCH=arm go build -o $SCRIPTDIR/nshptt_netbsd_arm $SCRIPTDIR/../tool.go  
env GOOS=netbsd GOARCH=arm64 go build -o $SCRIPTDIR/nshptt_netbsd_arm64 $SCRIPTDIR/../tool.go  

env GOOS=freebsd GOARCH=amd64 go build -o $SCRIPTDIR/nshptt_freebsd_amd64 $SCRIPTDIR/../tool.go  
env GOOS=freebsd GOARCH=386 go build -o $SCRIPTDIR/nshptt_freebsd_386 $SCRIPTDIR/../tool.go  
env GOOS=freebsd GOARCH=arm go build -o $SCRIPTDIR/nshptt_freebsd_arm $SCRIPTDIR/../tool.go  
env GOOS=freebsd GOARCH=arm64 go build -o $SCRIPTDIR/nshptt_freebsd_arm64 $SCRIPTDIR/../tool.go  