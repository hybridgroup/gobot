PACKAGES := gobot gobot/api gobot/platforms/firmata/client gobot/platforms/intel-iot/edison gobot/sysfs $(shell ls ./platforms | sed -e 's/^/gobot\/platforms\//')
.PHONY: test cover robeaux

test:
	for package in $(PACKAGES) ; do \
		go test -a github.com/hybridgroup/$$package ; \
	done ; \

cover:
	echo "mode: set" > profile.cov ; \
	for package in $(PACKAGES) ; do \
		go test -a -coverprofile=tmp.cov github.com/hybridgroup/$$package ; \
		cat tmp.cov | grep -v "mode: set" >> profile.cov ; \
	done ; \
	rm tmp.cov ; \

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
