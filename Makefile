# include also examples in other than ./examples folder
ALL_EXAMPLES := $(shell grep -l -r --include "*.go" 'build example' ./)
# prevent examples with joystick (sdl2) and gocv (opencv) dependencies
EXAMPLES := $(shell grep -L 'joystick' $$(grep -L 'gocv' $(ALL_EXAMPLES)))

.PHONY: test race cover robeaux test_with_coverage fmt_check examples $(EXAMPLES)

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

examples: $(EXAMPLES)

$(EXAMPLES):
	go build -o /tmp/gobot_examples/$@ ./$@
