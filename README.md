## Nullix Server Hardware Performance Test Tool (NSHPTT)
A simple tool to test raw performance of any server, pc or machine. It features CPU and Disk stress tests and can/should be used to compare multiple machines to find bottlenecks.

I created this because after a huge migration from one hosting company/location to another with tons of virtual machines, there was a lot of performance issues at the new location, which where incredible hard to nail down. At the end we found disk performance problems which only showed up clearly when reading files > 100MB... Long story short, this tool was developed to find this issues and have reproducable and comparable results.

![Slideshow](media/slideshow.gif?raw=true "NSHPTT")

## Features

- Single Core Performance Test. The test just calculates square-root of random floats as often as it can, for a set amount of time.
  - For servers it's better to test single core performance, as all cores would probably have the same specs anyway. With this method you can compare different infrastructures better then with multi-core tests.
- Disk Write-Read-Delete Test. Does write, than read, than delete files with random byte contents. Size, location and amount of file can be configured.
  - This test is good to test a lot of smaller files. You can simulate even higher IOPS when you lower filesize and increase number of files.
- Disk Read-Only Test. Does only read a pre-generated set of files with random byte contents. Size, location and amount of file can be configured.
  - This test should be used to test large files that don't change on disk. You have to create them with `--create-test-files` before the testrun. Disks and storage caches behave differently for random files that always change and files that are almost static.
- Test network impact by first running the default tests to have a base result. Then change the `--test-folder` to a path that is network attached.
- All tests, file sizes and file numbers are configurable with a `config.json` file which is generated with `--create-config`
- Results are logged in several ways: In the console, as a viewable `.html` with charts and as a CSV file to be machine readable.
- Works probably on every architecture and platform, even if we can't test em all
- Ready to run binaries in our releases. If u don't trust the binaries, just review and run to source `tool.go` file with Go itself.


## How it works
The tool run several tests, stressing different parts of your machine. The machine will be stressed for some time, so maybe don't run it in a production environment or test the config before on a testhost.
After the tests are stopped, a summary .html/.csv file is generated which you can view directly.

## Usage Pre-Built Binaries
Download a pre-built binary from "Releases" for your OS architecture and run it via command-line. The man page will show up that lists you all possible usages.

    nshptt_linux_amd64 --create-config
    nshptt_linux_amd64 --create-test-files
    nshptt_linux_amd64 --run

> Notice: The pre-compiled binaries can probably be detected by anti-virus software, as they are not signed and AV heuristics maybe struggle with that.
If you don't trust the pre-compiled binaries, feel free to directly use the `Usage directly with GO` variant. You can review and inspect the script before you run it.

## Usage directly with Golang

    go run tool.go --create-config
    go run tool.go --create-test-files
    go run tool.go --run


## Development

Always create an `issue` at github before you start changing things that you want to be merged into this repository.

- Requirements: [Download Go](https://go.dev/dl/) and [Install Go](https://go.dev/doc/install)
- If you do development on a windows machine, use WSL with a linux distro

## Build/Compile
- Requirements: run `bash setup.sh`

Just have a look at `build.sh` to see a list of all supported platforms and architectures and show to compile them with a single command line call.

Run `bash setup.sh` to build all into `build` directory.
