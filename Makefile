## Build lemming and deploy KNE topology
.PHONY: clean
clean:
	find integration_tests -name "topology.pb.txt" -exec kne delete {} \;; exit 0

.PHONY: load 
load:
	DOCKER_BUILDKIT=1 docker build . --target release -f Dockerfile.lemming -t "us-west1-docker.pkg.dev/openconfig-lemming/release/lemming:ga"
	kind load docker-image us-west1-docker.pkg.dev/openconfig-lemming/release/lemming:ga --name kne

.PHONY: buildfile
buildfile:
	go mod tidy
	bazel run //:gazelle -- update-repos -from_file=go.mod
	bazel run //:gazelle

.PHONY: load-debug
load-debug:
	DOCKER_BUILDKIT=1 docker build . --target debug -f Dockerfile.lemming -t "us-west1-docker.pkg.dev/openconfig-lemming/release/lemming:ga"
	kind load docker-image us-west1-docker.pkg.dev/openconfig-lemming/release/lemming:ga --name kne

## Run integration tests
.PHONY: itest
itest:
	bazel test //integration_tests/...

.PHONY: test
test:
	bazel test $(shell bazel query 'tests("//...") except "//integration_tests/..."')

.PHONY: test-race
test-race:
	# TODO: Fix race tests for lemming/gnmi and dataplane
	bazel test --@io_bazel_rules_go//go/config:race $(shell bazel query 'tests("//...") except "//integration_tests/..." except "//dataplane/..." except "//gnmi/..." ')
