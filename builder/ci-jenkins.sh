#!/bin/bash

# Jenkins Setting
# - Repository URL: git@github.com:QubitProducts/bamboo.git
# - Checkout to a sub-directory: go/src/github.com/QubitProducts/bamboo
# - Post build action / Archive the artifacts: *.deb
# - Test output: test_output/tests.xml

export BAMBOO_PROJECT_DIR=$WORKSPACE/go/src/github.com/QubitProducts/bamboo
export GOPATH=$WORKSPACE/go
export PATH=$GOPATH/bin:/usr/local/go/bin/:$PATH

# Cleanups
rm -f $WORKSPACE/*.deb
mkdir -p $WORKSPACE/test_output

# Install build and test dependencies
go get github.com/tools/godep
go get bitbucket.org/tebeka/go2xunit
go get -t github.com/smartystreets/goconvey


cd $BAMBOO_PROJECT_DIR
export _BAMBOO_VERSION=`(cat VERSION)`

go build
go test -v github.com/QubitProducts/bamboo/... | go2xunit > $WORKSPACE/test_output/tests.xml

# Requires fpm to build package
# Install fpm if missing
if ! gem spec fpm > /dev/null 2>&1; then gem install fpm; fi

# Build debian package
./builder/build.sh

# Copy files to workspace root directory
cp $BAMBOO_PROJECT_DIR/output/*.deb $WORKSPACE/
