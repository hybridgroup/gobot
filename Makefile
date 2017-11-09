.PHONY: test race cover robeaux examples deps test_with_coverage fmt_check

excluding_vendor := $(shell go list ./... | grep -v /vendor/)

# Run tests on all non-vendor directories
test:
	go test -v $(excluding_vendor)

# Run tests with race detection on all non-vendor directories
race:
	go test -race $(excluding_vendor)

# Check for code well-formedness
fmt_check:
	./ci/format.sh

# Test and generate coverage
test_with_coverage:
	./ci/test.sh

deps:
ifeq (,$(shell which dep))
	$(error dep tool not found! https://github.com/golang/dep is required to install Gobot deps)
endif
	dep ensure


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
