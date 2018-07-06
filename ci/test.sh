#!/bin/bash

## Only uncomment the below for debugging
#set -euxo pipefail

LOCAL_GO_VERSION=$(go version | awk -F' ' '{print $3}' | tr -d '[:space:]')
GO_VERSION="${TRAVIS_GO_VERSION:=$LOCAL_GO_VERSION}"
TIP_VERSION_IDENTIFIER="tip"
echo $GO_VERSION

# Hold the package names that contain failures
FAIL_PACKAGES=()

# OpenCV components to link in CGO compile
OPENCV_LDFLAGS="-lopencv_core -lopencv_face -lopencv_videoio -lopencv_imgproc -lopencv_highgui -lopencv_imgcodecs -lopencv_objdetect -lopencv_features2d -lopencv_video -lopencv_dnn -lopencv_xfeatures2d"

# Use $HOME on Travis
# Use /usr/local on local
if [[ $TRAVIS == "true" ]]; then
  export CGO_CPPFLAGS="-I${HOME}/usr/include"
  export CGO_LDFLAGS="-L${HOME}/usr/lib $OPENCV_LDFLAGS"
else
  export CGO_CPPFLAGS="-I/usr/local/include"
  export CGO_LDFLAGS="-L/usr/local/lib $OPENCV_LDFLAGS"
fi


pushd $PWD/..
	# Set up coverage report file
	COVERAGE_REPORT_LOCATION="./profile.cov"
	echo "" > $COVERAGE_REPORT_LOCATION

	# Exclude vendor etc.
	EXCLUDING_VENDOR=$(go list ./... | grep -Ev 'vendor|common|client|cli|examples|robeaux')

	# Iterate over all non-vendor packages and run tests with coverage
	for package in $EXCLUDING_VENDOR; do \
    if [ $GO_VERSION == $TIP_VERSION_IDENTIFIER ]; then
		# `go test` runs a vet subset as of this Change https://go-review.googlesource.com/c/go/+/74356
		  result=$(go test -vet=off -covermode=count -coverprofile=tmp.cov $package)
    else
		  result=$(go test -covermode=count -coverprofile=tmp.cov $package)
    fi
		if [ $? -ne 0 ]; then
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
