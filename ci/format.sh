#!/bin/bash

## Only uncomment the below for debugging
#set -euxo pipefail

pushd $PWD/..
  GO_FILES_EXCLUDING_VENDOR=$(find . -type f -name '*.go' -not -path "./vendor/*")
  FMT_RESULTS=$(gofmt -l $GO_FILES_EXCLUDING_VENDOR)
  FMT_RESULTS_COUNT=$(echo $FMT_RESULTS | wc -l) # returns one empty line when everything passes
  if [ "$FMT_RESULTS_COUNT" -gt 1 ]; then
    # some files have formatting errors
    echo "--- gofmt found errors found in the following files:"
    echo $FMT_RESULTS
    exit 1
  else
    # no files have formatting errors
    exit 0
  fi
popd
