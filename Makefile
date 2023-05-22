# include also examples in other than ./examples folder
ALL_EXAMPLES := $(shell grep -l -r --include "*.go" 'build example' ./)
# prevent examples with gocv (opencv) dependencies
EXAMPLES_NO_GOCV := $(shell grep -L 'gocv' $(ALL_EXAMPLES))
# prevent examples with joystick (sdl2) dependencies
EXAMPLES_NO_JOYSTICK := $(shell grep -L 'joystick' $(ALL_EXAMPLES))
# prevent examples with joystick (sdl2) and gocv (opencv) dependencies
EXAMPLES_NO_GOCV_JOYSTICK := $(shell grep -L 'joystick' $$(grep -L 'gocv' $(EXAMPLES_NO_GOCV)))
# used examples
EXAMPLES := $(EXAMPLES_NO_GOCV_JOYSTICK)

.PHONY: test test_race test_cover robeaux version_check fmt_check fmt_fix examples examples_check $(EXAMPLES)

# opencv platform currently skipped to prevent install of preconditions
including_except := $(shell go list ./... | grep -v platforms/opencv)

# Run tests on nearly all directories without test cache
test:
	cd ./v2 ; \
	go test -count=1 -v $(including_except)

# Run tests with race detection
test_race:
	cd ./v2 ; \
	go test -race $(including_except)

# Test, generate and show coverage in browser
test_cover:
	cd ./v2 ; \
	go test -v $(including_except) -coverprofile=coverage.txt ; \
	go tool cover -html=coverage.txt ; \

robeaux:
ifeq (,$(shell which go-bindata))
	$(error robeaux not built! https://github.com/jteeuwen/go-bindata is required to build robeaux assets )
endif
	cd ./v2 ; \
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
	cd ./v2 ; \
	@gv=$$(echo $$(go version) | sed "s/^.* go\([0-9].[0-9]*\).*/\1/") ; \
	mv=$$(grep -m 1 'go 1.' ./go.mod | sed "s/^go \([0-9].[0-9]*\).*/\1/") ; \
	echo "go: $${gv}.*, go.mod: $${mv}" ; \
	if [ "$${gv}" != "$${mv}" ]; then exit 50; fi ; \

# Check for bad code style and other issues
fmt_check:
	cd ./v2 ; \
	gofmt -l ./ ; \
	golangci-lint run -v

# Fix bad code style (will only be executed, on version match)
fmt_fix: version_check
	cd ./v2 ; \
	go fmt ./...

examples: $(EXAMPLES)

examples_check: 
	$(MAKE) CHECK=ON examples

$(EXAMPLES):
ifeq ($(CHECK),ON)
	go vet ./$@
else
	cd ./v2 ; \
	go build -o /tmp/gobot_examples/$@ ./$@
endif
