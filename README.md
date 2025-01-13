## Nullix Server Hardware Performance Test Tool (NSHPTT)
A simple tool to test raw performance of any server, pc or machine. It features CPU and Disk stress tests and can/should be used to compare multiple machines to find bottlenecks.

I created this because after a huge migration from one hosting company/location to another with tons of virtual machines, there was a lot of performance issues at the new location, which where incredible hard to nail down. At the end it was a disk performance problem which only showed up when reading files > 100MB... Long story short, this tool was developed to find this issues.

![Slideshow](media/slideshow.gif?raw=true "NSHPTT")

## How it works
The tool run several tests, stressing different parts of your machine. The machine will be stressed for some time, so maybe don't run it in a production environment or test the config before on a testhost.
After the tests are stopped, a summary .html/.csv file is generated which you can view directly.

## Usage Pre-Built Binaries
Download a pre-built binary from "Releases" for your OS architecture and run it via command-line. The man page will show up that lists you all possible usages.

    nshptt_linux_amd64 --create-config
    nshptt_linux_amd64 --create-test-files
    nshptt_linux_amd64 --run

> Notice: The pre-compiled binaries can be detected by anti-virus software, as they are not signed and read/write files from the testfolder.
If you don't trust the pre-compiled binaries, feel free to directly use the `Usage directly with GO` variant. You can review and inspect the script before you run it.

## Usage directly with Golang

    go run tool.go --create-config
    go run tool.go --create-test-files
    go run tool.go --run


## Development

Always create an `issue` at github before you start changing things that you want to be merged into this repository.

- Requirements: [Download Go](https://go.dev/dl/) and [Install Go](https://go.dev/doc/install)
- If you do development on a windows machine, use WSL with a linux distro

## Build
- Requirements: `build/setup.sh`

Just have a look at `build/build.sh` to see a list of all supported platforms and architectures and show to compile them with a single command line call.
