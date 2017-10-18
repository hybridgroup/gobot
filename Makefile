.PHONY: test race cover robeaux examples deps

test:
	go test ./...

race:
	go test ./... -race

cover:
	echo "" > profile.cov
	for package in $$(go list ./...) ; do \
		go test -covermode=count -coverprofile=tmp.cov $$package ; \
		if [ -f tmp.cov ]; then \
			cat tmp.cov >> profile.cov ; \
			rm tmp.cov ; \
		fi ; \
	done

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
	go get -d -v \
		github.com/bmizerany/pat \
		github.com/codegangsta/cli \
		github.com/currantlabs/ble \
		github.com/donovanhide/eventsource \
		github.com/eclipse/paho.mqtt.golang \
		github.com/hashicorp/go-multierror \
		github.com/hybridgroup/go-ardrone/client \
		gocv.io/x/gocv \
		github.com/mgutz/logxi/v1 \
		github.com/nats-io/nats \
		github.com/sigurn/crc8 \
		go.bug.st/serial.v1 \
		github.com/veandco/go-sdl2/sdl \
		golang.org/x/net/websocket \
		golang.org/x/exp/io/spi \
		golang.org/x/sys/unix
