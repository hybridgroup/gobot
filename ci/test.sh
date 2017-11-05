#!/bin/bash

## Only uncomment the below for debugging
#set -euxo pipefail

# Hold the package names that contain failures
FAIL_PACKAGES=()

pushd $PWD/..
	# Set up coverage report file
	COVERAGE_REPORT_LOCATION="./profile.cov"
	echo "" > $COVERAGE_REPORT_LOCATION

	# Exclude vendor in the same way as Makefile does
	EXCLUDING_VENDOR=$(go list ./... | grep -v /vendor/)

	# Iterate over all non-vendor packages and run tests with coverage
	for package in $EXCLUDING_VENDOR; do \
		result=$(go test -covermode=count -coverprofile=tmp.cov $package)
		# If a `go test` command has failures, it will exit 1
		# a go vet check should exit 2 (https://github.com/golang/go/blob/master/src/cmd/vet/all/main.go#L58)
		# `go test` contains `vet` was introduced in this Change https://go-review.googlesource.com/c/go/+/74356
		if [ $? -eq 1 ]; then
			FAIL_PACKAGES+=($package);
		fi;
			echo "$result"
		if [ -f tmp.cov ]; then
			cat tmp.cov >> profile.cov;
			rm tmp.cov;
		fi;
	done

	# exit 1 if there have been any test failures
	if [ ${#FAIL_PACKAGES[@]} -ne 0 ]; then
		exit 1
	fi;
popd
