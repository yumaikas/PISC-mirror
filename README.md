# PISC
Position Independent Source Code. A small, stack-based, concatenative language.

## About 

PISC's source and documentation are hosted at [https://pisc.junglecoder.com](https://pisc.junglecoder.com).
A quick starter can be found at [PISC in Y Minutes](https://pisc.junglecoder.com/home/apps/fossil/PISC.fossil/wiki?name=PISC+in+Y+Minutes)

## Building 

PISC requires go 1.7 or higher to build. Installation instructions [here](https://golang.org/doc/install)

Once go is installed, you can run `go get` to fetch the depenencies for PISC, and `go build -o pisc` to get a PISC executable. You can launch a REPL with `pisc -i` and play with it there.

## Flags

The PISC binary accepts 4 flags -i or --interactive, -c  or --command, -f or --file (useful for shebang lines), and -help, which displays their documentation
