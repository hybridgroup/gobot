# include also examples in other than ./examples folder
ALL_EXAMPLES := $(shell grep -l -r --include "*.go" 'build example' ./)
# prevent examples with gocv (opencv) dependencies
EXAMPLES_NO_GOCV := $(shell grep -L 'gocv' $(ALL_EXAMPLES))
# used examples
EXAMPLES := $(EXAMPLES_NO_GOCV)

.PHONY: test test_race test_cover robeaux version_check fmt_check fmt_fix examples examples_check examples_fmt_fix $(EXAMPLES)

# opencv platform currently skipped to prevent install of preconditions
including_except := $(shell go list ./... | grep -v platforms/opencv)

# Run tests on nearly all directories without test cache, with race detection
test_race:
	go test -failfast -count=1 -race $(including_except) -tags libusb

# Run tests on nearly all directories without test cache
test:
	go test -failfast -count=1 $(including_except)

# Test, generate and show coverage in browser
test_cover:
	go test -v $(including_except) -coverprofile=coverage.txt ; \
	go tool cover -html=coverage.txt ; \

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

# Check for installed and module version match. Will exit with code 50 if not match.
# There is nothing bad in general, if you program with a higher version.
# At least the recipe "fmt_fix" will not work in that case.
version_check:
	@gv=$$(echo $$(go version) | sed "s/^.* go\([0-9].[0-9]*\).*/\1/") ; \
	mv=$$(grep -m 1 'go 1.' ./go.mod | sed "s/^go \([0-9].[0-9]*\).*/\1/") ; \
	echo "go: $${gv}.*, go.mod: $${mv}" ; \
	if [ "$${gv}" != "$${mv}" ]; then exit 50; fi ; \

# Check for bad code style and other issues (gofumpt and gofmt check is activated for the linter)
fmt_check:
	golangci-lint run -v

# Fix bad code style (the go version will be automatically obtained from go.mod)
fmt_fix: 
	$(MAKE) version_check || true
	gofumpt -l -w .
	golangci-lint run -v --fix

examples: $(EXAMPLES)

examples_check: 
	$(MAKE) CHECK=ON examples

examples_fmt_fix: 
	$(MAKE) CHECK=FMT examples

$(EXAMPLES):
ifeq ($(CHECK),ON)
	go vet -tags libusb ./$@
else ifeq ($(CHECK),FMT)
	gofumpt -l -w ./$@
	golangci-lint run ./$@ --fix --build-tags example,libusb --disable forcetypeassert --disable noctx
else
	go build -tags libusb -o /tmp/gobot_examples/$@ ./$@
endif
