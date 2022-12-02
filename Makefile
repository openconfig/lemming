## Build lemming and deploy KNE topology
.PHONY: clean
clean:
	find integration_tests -name "topology.pb.txt" -exec kne delete {} \;; exit 0

.PHONY: load 
load:
	DOCKER_BUILDKIT=1 docker build . -f Dockerfile.lemming -t "lemming:latest"
	kind load docker-image lemming:latest --name kne


## Run integration tests
.PHONY: itest
itest:
	go test -v ./integration_tests/...

.PHONY: test
test:
	go test $(shell go list ./... | grep -v integration_test)
