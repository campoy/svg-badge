#! /bin/bash

set -e

tmp=$(mktemp -d)
dest="$tmp/src/github.com/campoy/svg-badge"
mkdir -p $dest

cp -r vendor/* $tmp/src
for f in $(ls | grep -v vendor); do
    cp -r $f $tmp/src/github.com/campoy/svg-badge
done

GOPATH=$tmp gcloud app deploy $dest/gae/app.yaml