.PHONY: test race cover robeaux examples test_with_coverage fmt_check

# opencv platform currently skipped to prevent install of preconditions
including_except := $(shell go list ./... | grep -v platforms/opencv)

# Run tests on nearly all directories
test:
	go test -v $(including_except)

# Run tests with race detection
race:
	go test -race $(including_except)

# Check for code well-formedness
fmt_check:
	./ci/format.sh

# Test and generate coverage
test_with_coverage:
	./ci/test.sh

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
