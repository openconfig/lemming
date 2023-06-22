## Build lemming and deploy KNE topology
.PHONY: clean
clean:
	find integration_tests -name "topology.pb.txt" -exec kne delete {} \;; exit 0

.PHONY: load-operator 
load-operator:
	bazel build //operator:image-tar
	kind load image-archive bazel-bin/operator/image-tar/tarball.tar --name kne

.PHONY: load 
load:
	bazel build //cmd/lemming:image-tar
	kind load image-archive bazel-bin/cmd/lemming/image-tar/tarball.tar --name kne

.PHONY: buildfile
buildfile:
	go mod tidy
	bazel run //:gazelle -- update-repos -to_macro=repositories.bzl%go_repositories -from_file=go.mod
	bazel run //:gazelle

.PHONY: load-debug
load-debug:
	DOCKER_BUILDKIT=1 docker build . --target debug -f Dockerfile.lemming -t "us-west1-docker.pkg.dev/openconfig-lemming/release/lemming:ga"
	kind load docker-image us-west1-docker.pkg.dev/openconfig-lemming/release/lemming:ga --name kne

## Run integration tests
.PHONY: itest
itest:
	bazel test --test_output=errors --cache_test_results=no //integration_tests/...

.PHONY: test
test:
	bazel test --test_output=errors $(shell bazel query 'tests("//...") except "//integration_tests/..." except "//bgp/tests/local_tests/..."')
	bazel test --test_output=all --test_arg=-v=0 --test_arg=-alsologtostderr //bgp/tests/local_tests/...
	bazel test --test_output=all --test_arg=-v=0 --test_arg=-alsologtostderr //bgp/tests/local_tests/...

.PHONY: test-race
test-race:
	# TODO: Fix race tests for lemming/gnmi and dataplane
	# Failure in local_tests are due to GoBGP itself unable to issue a Stop
	# command without conflicting with the running server in another
	# thread.(e.g. TestRoutePropagation)
	bazel test --@io_bazel_rules_go//go/config:race --test_output=errors $(shell bazel query 'tests("//...") except "//integration_tests/..." except "//dataplane/..." except "//gnmi/..." except "//bgp/tests/local_tests/..."')
