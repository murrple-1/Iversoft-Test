#!/bin/sh

DIR=$(dirname $0)

cd $DIR/../

go get

go build -o iversoft.out

