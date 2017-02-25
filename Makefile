PACKAGES := gobot gobot/api gobot/drivers/gpio gobot/drivers/aio gobot/drivers/i2c gobot/platforms/firmata/client gobot/platforms/intel-iot/edison gobot/platforms/intel-iot/joule gobot/platforms/parrot/ardrone gobot/platforms/parrot/bebop gobot/platforms/parrot/minidrone gobot/platforms/sphero/ollie gobot/platforms/sphero/bb8 gobot/sysfs $(shell ls ./platforms | sed -e 's/^/gobot\/platforms\//')
.PHONY: test cover robeaux examples

test:
	go test -i ./...
	for package in $(PACKAGES) ; do \
		go test gobot.io/x/$$package ; \
	done ; \

cover:
	echo "" > profile.cov
	go test -i ./...
	for package in $(PACKAGES) ; do \
		go test -covermode=count -coverprofile=tmp.cov gobot.io/x/$$package ; \
		cat tmp.cov >> profile.cov ; \
		rm tmp.cov ; \
	done ; \

robeaux:
ifeq (,$(shell which go-bindata))
	$(error robeaux not built! https://github.com/jteeuwen/go-bindata is required to build robeaux assets )
endif
	cd api ; \
	npm install robeaux ; \
	cp -r node_modules/robeaux robeaux-tmp ; \
	cd robeaux-tmp ; \
	rm Makefile package.json README.markdown ; \
	touch css/fonts.css ; \
	echo "Updating robeaux..." ; \
	go-bindata -pkg="robeaux" -o robeaux.go -ignore=\\.git ./... ; \
	mv robeaux.go ../robeaux ; \
	cd .. ; \
	rm -rf robeaux-tmp/ ; \
	rm -rf node_modules/ ; \
	go fmt ./robeaux/robeaux.go ; \

EXAMPLES := $(shell ls examples/*.go | sed -e 's/examples\///')

examples:
	for example in $(EXAMPLES) ; do \
		go build -o /tmp/$$example examples/$$example ; \
	done ; \

deps:
	go get -d -v github.com/bmizerany/pat
	go get -d -v github.com/hybridgroup/go-ardrone/client
	go get -d -v github.com/mgutz/logxi/v1
	go get -d -v golang.org/x/sys/unix
	go get -d -v github.com/currantlabs/ble
	go get -d -v github.com/tarm/serial
	go get -d -v github.com/veandco/go-sdl2/sdl
	go get -d -v golang.org/x/net/websocket
	go get -d -v github.com/eclipse/paho.mqtt.golang
	go get -d -v github.com/nats-io/nats
	go get -d -v github.com/lazywei/go-opencv
	go get -d -v github.com/donovanhide/eventsource
	go get -d -v github.com/hashicorp/go-multierror
	go get -d -v github.com/sigurn/crc8
	go get -d -v github.com/codegangsta/cli
