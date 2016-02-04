#!/bin/bash
PACKAGES=('gobot' 'gobot/api' 'gobot/platforms/firmata/client' 'gobot/platforms/intel-iot/edison' 'gobot/sysfs' $(ls ./platforms | sed -e 's/^/gobot\/platforms\//'))
EXITCODE=0

echo "mode: set" > profile.cov
touch tmp.cov
for package in "${PACKAGES[@]}"
do
  go test -a -coverprofile=tmp.cov github.com/hybridgroup/$package
  if [ $? -ne 0 ]
  then
    EXITCODE=1
  fi
  cat tmp.cov | grep -v "mode: set" >> profile.cov
done

if [ $EXITCODE -ne 0 ]
then
  exit $EXITCODE
fi

export PATH=$PATH:$HOME/gopath/bin/
goveralls -v -coverprofile=profile.cov -service=travis-ci -repotoken=sFrR9ZmLP5FLc34lOaqir67RPzYOvFPUB
