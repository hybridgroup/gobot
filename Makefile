PACKAGES := gobot gobot/api $(shell ls ./platforms | sed -e 's/^/gobot\/platforms\//')

.PHONY: test cover robeaux 

test:
	for package in $(PACKAGES) ; do \
		go test github.com/hybridgroup/$$package ; \
	done ; \
 
cover:
	echo "mode: count" > profile.cov ; \
	for package in $(PACKAGES) ; do \
		go test -covermode=count -coverprofile=tmp.cov github.com/hybridgroup/$$package ; \
		cat tmp.cov | grep -v "mode: count" >> profile.cov ; \
	done ; \
	rm tmp.cov ; \

robeaux:
ifeq (,$(shell which go-bindata))
	$(error robeaux not built! https://github.com/jteeuwen/go-bindata is required to build robeaux assets )
endif
	cd api ; \
	git clone --depth 1 git://github.com/hybridgroup/robeaux.git ; \
	cd robeaux ; \
	rm fonts/* ; \
	rm Makefile package.json README.markdown robeaux.gemspec css/fonts.css ; \
	touch css/fonts.css ; \
	echo "Updating robeaux to $(shell git rev-parse HEAD)" ; \
	go-bindata -pkg="api" -o robeaux.go -ignore=\\.git ./... ; \
	mv robeaux.go .. ; \
	cd .. ; \
	rm -rf robeaux/ ; \
	go fmt ./robeaux.go ; \
	
