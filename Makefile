## Build lemming and deploy KNE topology
.PHONY: clean
clean:
	kne delete kne/twodut_topo.pb.txt; kne delete kne/topo.pb.txt; kne delete kne/twodut_oneotg_topo.pb.txt; exit 0

.PHONY: deploy 
deploy:
	kne/deploy.sh topo.pb.txt

.PHONY: deploy2
deploy2:
	kne/deploy.sh twodut_topo.pb.txt

.PHONY: deploy3
deploy3:
	kne/deploy.sh twodut_oneotg_topo.pb.txt

## Run integration tests
.PHONY: itest
itest:
	kne/setup.sh topo.pb.txt
	go test -v ./integration_tests -args -config $(shell pwd)/kne/config.yaml -testbed $(shell pwd)/kne/testbed.pb.txt

.PHONY: itest2
itest2:
	kne/setup.sh twodut_topo.pb.txt
	go test -v ./integration_tests/twodut_tests -args -config $(shell pwd)/kne/config.yaml -testbed $(shell pwd)/kne/twodut_testbed.pb.txt

.PHONY: itest3
itest3:
	kne/setup.sh twodut_oneotg_topo.pb.txt
	go test -v ./integration_tests/twodut_oneotg_tests -args -config $(shell pwd)/kne/config.yaml -testbed $(shell pwd)/kne/twodut_oneotg_testbed.pb.txt

.PHONY: test
test:
	go test $(shell go list ./... | grep -v integration_test)
