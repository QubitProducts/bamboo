export GOPATH=$WORKSPACE/go
export PATH=$GOPATH/bin:/usr/local/go/bin/:$PATH
rm -f $WORKSPACE/*.deb

mkdir -p $WORKSPACE/test_output

go get github.com/QubitProducts/bamboo
go get github.com/tools/godep
go get bitbucket.org/tebeka/go2xunit
go get -t github.com/smartystreets/goconvey

export BAMBOO_PROJECT_DIR=$WORKSPACE/go/src/github.com/QubitProducts/bamboo

cd $BAMBOO_PROJECT_DIR

godep restore
go build
go test -v bamboo/... | go2xunit > $WORKSPACE/test_output/tests.xml

# Requires fpm to build package
# Install fpm if missing
if ! gem spec fpm > /dev/null 2>&1; then gem install fpm; fi

# Build debian package
./builder/build.sh

# Copy files to workspace root directory
cp $BAMBOO_PROJECT_DIR/builder/*.deb $WORKSPACE/