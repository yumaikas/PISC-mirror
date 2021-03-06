# PISC
Position Independent Source Code. A small, stack-based, concatenative language.

## About 

PISC's source and documentation are hosted at [https://pisc.junglecoder.com/ ](https://pisc.junglecoder.com).
A quick starter can be found at [PISC in Y Minutes](https://pisc.junglecoder.com/home/apps/fossil/PISC.fossil/wiki?name=PISC+in+Y+Minutes)

## Building 

PISC requires go 1.9+ to build. Installation instructions [here](https://golang.org/doc/install)

Once go is installed, you'll (currently) need to run `git clone https://github.com/yumaikas/PISC-mirror "$GOPATH/src/pisc"`
Running `cd $GOPATH/src/pisc && go get -u && go build -o pisc` will fetch depenencies and get you a PISC executable. You can launch a REPL with `pisc -i` and play with it there.

## Playground

PISC has a playground at https://pisc.junglecoder.com/playground/

## Flags

Running `pisc help` will dump out current information about avaiable flags

