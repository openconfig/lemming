## Build lemming and deploy KNE topology
.PHONY: clean
clean:
	find integration_tests -name "topology.pb.txt" -exec kne delete {} \;; exit 0

.PHONY: load 
load:
	DOCKER_BUILDKIT=1 docker build . -f Dockerfile.lemming -t "us-west1-docker.pkg.dev/openconfig-lemming/release/lemming:ga"
	kind load docker-image us-west1-docker.pkg.dev/openconfig-lemming/release/lemming:ga --name kne


## Run integration tests
.PHONY: itest
itest:
	go list ./integration_tests/... |  while read -r test ; do echo "Running test $$test"; go test -v -count 1 "$$test"; done

.PHONY: test
test:
	go test $(shell go list ./... | grep -v integration_test)
