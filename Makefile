PACKAGES := "github.com/hybridgroup/gobot" "github.com/hybridgroup/gobot/api" "github.com/hybridgroup/gobot/platforms/ardrone" "github.com/hybridgroup/gobot/platforms/beaglebone" "github.com/hybridgroup/gobot/platforms/digispark" "github.com/hybridgroup/gobot/platforms/firmata" "github.com/hybridgroup/gobot/platforms/gpio" "github.com/hybridgroup/gobot/platforms/i2c" "github.com/hybridgroup/gobot/platforms/joystick" "github.com/hybridgroup/gobot/platforms/leap" "github.com/hybridgroup/gobot/platforms/neurosky" "github.com/hybridgroup/gobot/platforms/opencv" "github.com/hybridgroup/gobot/platforms/pebble" "github.com/hybridgroup/gobot/platforms/spark" "github.com/hybridgroup/gobot/platforms/sphero" 

test:
	go test ./...
 
cover:
	echo "mode: count" > profile.cov ; \
	for package in $(PACKAGES) ; do \
		go test -covermode=count -coverprofile=tmp.cov $$package ; \
		cat tmp.cov | grep -v "mode: count" >> profile.cov ; \
	done ; \
	rm tmp.cov ; \
