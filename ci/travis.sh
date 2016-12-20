#!/bin/bash
PACKAGES=('gobot' 'gobot/api' 'gobot/sysfs' 'gobot/drivers/gpio' 'gobot/drivers/aio' 'gobot/drivers/i2c' 'gobot/platforms/firmata/client' 'gobot/platforms/intel-iot/edison' 'gobot/platforms/intel-iot/joule' 'gobot/platforms/parrot/ardrone' 'gobot/platforms/parrot/bebop' 'gobot/platforms/parrot/minidrone' 'gobot/platforms/sphero/ollie' 'gobot/platforms/sphero/bb8' $(ls ./platforms | sed -e 's/^/gobot\/platforms\//'))
EXITCODE=0

echo "mode: set" > profile.cov
touch tmp.cov
for package in "${PACKAGES[@]}"
do
  go test -a -coverprofile=tmp.cov gobot.io/x/$package
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
