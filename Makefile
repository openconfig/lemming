## Build lemming and deploy KNE topology
.PHONY: clean
clean:
	kne delete kne/twodut_topo.pb.txt; kne delete kne/topo.pb.txt

.PHONY: deploy 
deploy:
	kne/deploy.sh topo.pb.txt

.PHONY: deploy2
deploy2:
	kne/deploy.sh twodut_topo.pb.txt

## Run integration tests
.PHONY: itest
itest:
	kne/setup.sh topo.pb.txt
	go test -v ./integration_tests -args -config $(shell pwd)/kne/config.yaml -testbed $(shell pwd)/kne/testbed.pb.txt

.PHONY: itest2
itest2:
	kne/setup.sh twodut_topo.pb.txt
	go test -v ./integration_tests/twodut_tests -args -config $(shell pwd)/kne/config.yaml -testbed $(shell pwd)/kne/twodut_testbed.pb.txt

.PHONY: test
test:
	go test $(shell go list ./... | grep -v integration_test)
