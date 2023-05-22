# this file is just a forwarder to the folder with go.mod for common use cases
# it is working since Go 1.18 is installed locally

gomoddir := $(shell go list -f '{{.Dir}}' -m)


.PHONY: test fmt_check examples_check

test:
	cd $(gomoddir) && make test && cd ..

fmt_check:
	cd $(gomoddir) && make fmt_check && cd ..

examples_check:
	cd $(gomoddir) && make examples_check && cd ..