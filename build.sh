#!/usr/bin/bash
SCRIPTDIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
BUILDDIR=$SCRIPTDIR/build

rm -f $BUILDDIR/nshptt_*
rm -f $BUILDDIR/rsrc_*.syso

$BUILDDIR/go-winres make
mv $SCRIPTDIR/rsrc_*.syso $BUILDDIR

env GOOS=aix GOARCH=ppc64 go build -o $BUILDDIR/nshptt_aix_ppc64 $BUILDDIR/../tool.go 

# env GOOS=android GOARCH=amd64 go build -o $BUILDDIR/nshptt_android_amd64.exe $BUILDDIR/../tool.go  
# env GOOS=android GOARCH=386 go build -o $BUILDDIR/nshptt_android_386.exe $BUILDDIR/../tool.go  
# env GOOS=android GOARCH=arm go build -o $BUILDDIR/nshptt_android_arm.exe $BUILDDIR/../tool.go  
env GOOS=android GOARCH=arm64 go build -o $BUILDDIR/nshptt_android_arm64 $BUILDDIR/../tool.go  

env GOOS=darwin GOARCH=amd64 go build -o $BUILDDIR/nshptt_darwin_amd64 $BUILDDIR/../tool.go  
env GOOS=darwin GOARCH=arm64 go build -o $BUILDDIR/nshptt_darwin_arm64 $BUILDDIR/../tool.go  

# env GOOS=ios GOARCH=amd64 go build -o $BUILDDIR/nshptt_ios_amd64.exe $BUILDDIR/../tool.go  
# env GOOS=ios GOARCH=arm64 go build -o $BUILDDIR/nshptt_ios_arm64.exe $BUILDDIR/../tool.go  

env GOOS=windows GOARCH=amd64 go build -o $BUILDDIR/nshptt_win_amd64.exe $BUILDDIR/../tool.go && $BUILDDIR/go-winres patch $BUILDDIR/nshptt_win_amd64.exe
env GOOS=windows GOARCH=386 go build -o $BUILDDIR/nshptt_win_386.exe $BUILDDIR/../tool.go  && $BUILDDIR/go-winres patch $BUILDDIR/nshptt_win_386.exe
env GOOS=windows GOARCH=arm go build -o $BUILDDIR/nshptt_win_arm.exe $BUILDDIR/../tool.go  
env GOOS=windows GOARCH=arm64 go build -o $BUILDDIR/nshptt_win_arm64.exe $BUILDDIR/../tool.go  
rm -f $BUILDDIR/nshptt_win_*.bak

env GOOS=linux GOARCH=amd64 go build -o $BUILDDIR/nshptt_linux_amd64 $BUILDDIR/../tool.go  
env GOOS=linux GOARCH=386 go build -o $BUILDDIR/nshptt_linux_386 $BUILDDIR/../tool.go  
env GOOS=linux GOARCH=arm go build -o $BUILDDIR/nshptt_linux_arm $BUILDDIR/../tool.go  
env GOOS=linux GOARCH=arm64 go build -o $BUILDDIR/nshptt_linux_arm64 $BUILDDIR/../tool.go  

env GOOS=openbsd GOARCH=amd64 go build -o $BUILDDIR/nshptt_openbsd_amd64 $BUILDDIR/../tool.go  
env GOOS=openbsd GOARCH=386 go build -o $BUILDDIR/nshptt_openbsd_386 $BUILDDIR/../tool.go  
env GOOS=openbsd GOARCH=arm go build -o $BUILDDIR/nshptt_openbsd_arm $BUILDDIR/../tool.go  
env GOOS=openbsd GOARCH=arm64 go build -o $BUILDDIR/nshptt_openbsd_arm64 $BUILDDIR/../tool.go  

env GOOS=netbsd GOARCH=amd64 go build -o $BUILDDIR/nshptt_netbsd_amd64 $BUILDDIR/../tool.go  
env GOOS=netbsd GOARCH=386 go build -o $BUILDDIR/nshptt_netbsd_386 $BUILDDIR/../tool.go  
env GOOS=netbsd GOARCH=arm go build -o $BUILDDIR/nshptt_netbsd_arm $BUILDDIR/../tool.go  
env GOOS=netbsd GOARCH=arm64 go build -o $BUILDDIR/nshptt_netbsd_arm64 $BUILDDIR/../tool.go  

env GOOS=freebsd GOARCH=amd64 go build -o $BUILDDIR/nshptt_freebsd_amd64 $BUILDDIR/../tool.go  
env GOOS=freebsd GOARCH=386 go build -o $BUILDDIR/nshptt_freebsd_386 $BUILDDIR/../tool.go  
env GOOS=freebsd GOARCH=arm go build -o $BUILDDIR/nshptt_freebsd_arm $BUILDDIR/../tool.go  
env GOOS=freebsd GOARCH=arm64 go build -o $BUILDDIR/nshptt_freebsd_arm64 $BUILDDIR/../tool.go  