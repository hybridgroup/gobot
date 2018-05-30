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
	# Set up core coverage report file
	COVERAGE_REPORT_LOCATION="./core.cov"
	echo "" > $COVERAGE_REPORT_LOCATION

	result=$(go test -vet=off -covermode=count -coverprofile=core.cov ./gobot/. ./gobot/api/ ./gobot/sysfs/)
	if [ $? -ne 0 ]; then
		FAIL_PACKAGES+="CORE";
	fi;
	echo "$result"

	# Set up platforms coverage report file
	PLATFORMS_COVERAGE_REPORT_LOCATION="./platforms.cov"
	echo "" > $PLATFORMS_COVERAGE_REPORT_LOCATION

	result=$(go test -vet=off -covermode=count -coverprofile=platforms.cov ./gobot/platforms/...)
	if [ $? -ne 0 ]; then
		FAIL_PACKAGES+="PLATFORMS";
	fi;
	echo "$result"

	# Set up drivers coverage report file
	DRIVERS_COVERAGE_REPORT_LOCATION="./drivers.cov"
	echo "" > $DRIVERS_COVERAGE_REPORT_LOCATION

	result=$(go test -vet=off -covermode=count -coverprofile=platforms.cov ./gobot/drivers/...)
	if [ $? -ne 0 ]; then
		FAIL_PACKAGES+="PLATFORMS";
	fi;
	echo "$result"

	# exit 1 if there have been any test failures
	if [ ${#FAIL_PACKAGES[@]} -ne 0 ]; then
		exit 1
	fi;
popd
