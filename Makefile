## Build lemming and deploy KNE topology
.PHONY: deploy
deploy:
	kne/deploy.sh

## Run integration tests
.PHONY: itest
itest:
	go test -v ./integration_tests -args -config $(shell pwd)/kne/config.yaml -testbed $(shell pwd)/kne/testbed.pb.txt
