#!/bin/bash

# Delete every *.go files except hack.go file
function safeDelete() {
    mv ./hack.go ./hack 2> /dev/null
    rm -rf ./*.go 2> /dev/null
    mv ./hack ./hack.go 2> /dev/null
}

###

# git clone https://chromium.googlesource.com/webm/libvpx libvpxsrc
# git clone git@github.com:gotranspile/cxgo.git cxgosrc

sed -i 's+#define fread wrap_fread++g' ./libvpxsrc/tools_common.c
sed -i 's+ INT_MIN+ -2147483648+g' ./libvpxsrc/args.c
sed -i 's+ INT_MAX+ 2147483647+g' ./libvpxsrc/args.c
sed -i 's+ UINT_MAX+ 4294967295+g' ./libvpxsrc/args.c
sed -i 's+#include "string.h"+#include <string.h>+g' ./libvpxsrc/vp8/common/quant_common.h
cp -rf  ./includes ./libvpxsrc
cp -rf  ./include ./libvpxsrc

# https://unix.stackexchange.com/questions/77127/rm-rf-all-files-and-all-hidden-files-without-error
#rm -rf ./vpx/{*,.*}/ 2> /dev/null


# compile internal package
cd ./internal  || exit
safeDelete
cxgo -c internal.yml
rm go.mod

# compile vpx package
cd ./vpx  || exit
safeDelete
cxgo -c vpx.yml
rm go.mod
cd ../

# compile scale package
cd ./scale  || exit
safeDelete
cxgo -c scale.yml
rm go.mod
cd ../

# compile ports package
cd ./ports  || exit
safeDelete
cxgo -c ports.yml
rm go.mod
cd ../


# compile mem package
cd ./mem  || exit
safeDelete
cxgo -c mem.yml
rm go.mod
cd ../

# compile dsp package
cd ./dsp  || exit
safeDelete
cxgo -c dsp.yml
rm go.mod
cd ../

# compile util package
cd ./util  || exit
safeDelete
cxgo -c util.yml
rm go.mod
cd ../


# compile vp8 package
cd ./vp8  || exit
safeDelete
cxgo -c vp8.yml
rm go.mod
cd ../

# compile vp9 package
cd ./vp9  || exit
safeDelete
cxgo -c vp9.yml
rm go.mod
cd ../

# compile vpxdecgo package
cd ../
cxgo -c vpxdecgo.yml
rm ./cmd/vpxdecgo/go.mod
go mod tidy
go fmt ./...