PACKAGES := gobot gobot/api $(shell ls ./platforms | sed -e 's/^/gobot\/platforms\//')

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
