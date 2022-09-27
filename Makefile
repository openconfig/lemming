## Build lemming and deploy KNE topology
.PHONY: clean
clean:
	kne delete kne/twodut_topo.pb.txt; kne delete kne/topo.pb.txt

.PHONY: deploy
deploy:
	kne/deploy.sh

.PHONY: deploy2
deploy2:
	kne/deploytwodut.sh

## Run integration tests
.PHONY: itest
itest:
	go test -v ./integration_tests -args -config $(shell pwd)/kne/config.yaml -testbed $(shell pwd)/kne/testbed.pb.txt

.PHONY: itest2
itest2:
	go test -v ./integration_tests/twodut_tests -args -config $(shell pwd)/kne/twodut_config.yaml -testbed $(shell pwd)/kne/twodut_testbed.pb.txt
