#!/bin/bash
set -euo pipefail

# Print commands as executed
#set -x


# Hold the package names that contain failures
fail_packages=()

pushd $PWD/..
	# Set up coverage report file
	COVERAGE_REPORT_LOCATION="./profile.cov"
	echo "" > $COVERAGE_REPORT_LOCATION

	# Exclude vendor in the same way as Makefile does
	EXCLUDING_VENDOR=$(go list ./... | grep -v /vendor/)

	# Iterate over all non-vendor packages and run tests with coverage
	for package in $EXCLUDING_VENDOR; do \
		result=$(go test -covermode=count -coverprofile=tmp.cov $package)
		echo $result

		# If a `go test` command has failures, it will exit 1
		# a go vet check should exit 2
		if [ $? -eq 1 ]; then
			fail_packages+=($package);
		fi;
		if [ -f tmp.cov ]; then
			cat tmp.cov >> profile.cov;
			rm tmp.cov;
		fi;
	done

	# exit 1 if there have been any test failures
	if [ ${#fails[@]} -ne 0 ]; then
		exit 1
	fi;
popd
