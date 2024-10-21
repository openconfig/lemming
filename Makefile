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
	docker load -i bazel-bin/cmd/lemming/image-tar/tarball.tar
	kind load docker-image us-west1-docker.pkg.dev/openconfig-lemming/release/lemming:ga --name kne

.PHONY: buildfile
buildfile:
	go mod tidy
	bazel run //:gazelle -- update-repos -to_macro=repositories.bzl%go_repositories -from_file=go.mod -prune
	bazel run //:gazelle

.PHONY: genprotos
genprotos:
	tools/genproto.sh

.PHONY: load-debug
load-debug:
	DOCKER_BUILDKIT=1 docker build . --target debug -f Dockerfile.lemming -t "us-west1-docker.pkg.dev/openconfig-lemming/release/lemming:ga"
	kind load docker-image us-west1-docker.pkg.dev/openconfig-lemming/release/lemming:ga --name kne

## Run integration tests
## TODO: Reenable BGP triggered GUE tests once it works.
.PHONY: itest
itest:
	bazel test --flaky_test_attempts=3 --test_output=errors --cache_test_results=no $(shell bazel query 'tests("//...") except (//integration_tests/twodut_oneotg_tests/bgp_triggered_gue:bgp_triggered_gue_test + attr(size, small, tests("//...")) + attr(size, medium, tests("//..."))) ') 

.PHONY: test
test:
	bazel test --test_output=errors $(shell bazel query 'attr(size, small, tests("//...")) +  attr(size, medium, tests("//..."))')

.PHONY: coverage
coverage:
	bazel coverage --test_output=errors --combined_report=lcov  $(shell bazel query 'attr(size, small, tests("//...")) + attr(size, medium, tests("//..."))')

.PHONY: test-race
test-race:
	# TODO: Fix race tests for lemming/gnmi and dataplane
	# Failure in local_tests are due to GoBGP itself unable to issue a Stop
	# command without conflicting with the running server in another
	# thread.(e.g. TestRoutePropagation)
	bazel test --@io_bazel_rules_go//go/config:race --test_output=errors $(shell bazel query 'tests("//...") except "//integration_tests/..." except "//dataplane/..." except "//gnmi/..." except "//bgp/tests/local_tests/..."')

PROTOS = $(wildcard dataplane/proto/sai/*.proto)
PROTO_SRC = $(patsubst dataplane/proto/sai/%.proto, dataplane/proto/sai/%.pb.cc, $(PROTOS))
PROTO_GRPC_SRC = $(patsubst dataplane/proto/sai/%.proto, dataplane/proto/sai/%.grpc.pb.cc, $(PROTOS))
PROTO_OBJ = $(patsubst dataplane/proto/sai/%.proto, dataplane/proto/sai/%.pb.o, $(PROTOS))
GRPC_OBJ = $(patsubst dataplane/proto/sai/%.proto, dataplane/proto/sai/%.grpc.pb.o, $(PROTOS))
SAI_SRC = $(wildcard dataplane/standalone/sai/*.cc)
SAI_OBJ = $(patsubst dataplane/standalone/sai/%.cc, dataplane/standalone/sai/%.o, $(SAI_SRC))

.PHONY: sai-clean
sai-clean:
	rm dataplane/proto/sai/*.cc dataplane/proto/sai/*.h dataplane/proto/sai/*.o
	rm dataplane/standalone/sai/*.o
	rm -rf dataplane/standalone/packetio/packetio.a dataplane/standalone/packetio/packetio.h libsai.so pkg

$(PROTO_SRC) $(PROTO_GRPC_SRC) &:
	protoc dataplane/proto/sai/*.proto --cpp_out=. --grpc_out=. --plugin=protoc-gen-grpc=`which grpc_cpp_plugin` --experimental_allow_proto3_optional
    
dataplane/proto/sai/%.pb.o: dataplane/proto/sai/%.pb.cc
	g++ -fPIC -c $< -I . -o $@

dataplane/proto/sai/%.grpc.pb.o:  dataplane/proto/sai/%.grpc.pb.cc
	g++ -fPIC -c $< -I . -o $@

dataplane/standalone/sai/%.o: dataplane/standalone/sai/%.cc $(PROTO_SRC)
	g++ -fPIC -c $< -o $@ -I . -I external/com_github_opencomputeproject_sai -I external/com_github_opencomputeproject_sai/inc  -I external/com_github_opencomputeproject_sai/experimental


libsai.so: $(PROTO_OBJ) $(GRPC_OBJ) $(SAI_OBJ)
	g++ -fPIC -o libsai.so -shared dataplane/standalone/entrypoint.cc dataplane/proto/sai/*.o dataplane/standalone/sai/*.o -lglog -lprotobuf -lgrpc++ -I . -I external/com_github_opencomputeproject_sai -I external/com_github_opencomputeproject_sai/inc -I external/com_github_opencomputeproject_sai/experimental

define DEB_CONTROL =
Package: lucius-libsai
Version: 1.0-3
Maintainer: OpenConfig
Architecture: amd64
Description: SAI implementation for lucius
endef
export DEB_CONTROL

lucius-libsai: libsai.so 
	rm -rf pkg/
	mkdir -p pkg/lucius-libsai/DEBIAN
	chmod 0755 pkg/lucius-libsai/DEBIAN
	mkdir -p pkg/lucius-libsai/usr/lib/x86_64-linux-gnu
	cp libsai.so pkg/lucius-libsai/usr/lib/x86_64-linux-gnu/libsaivs.so.0.0.0
	cd pkg/lucius-libsai/usr/lib/x86_64-linux-gnu && ln -s libsaivs.so.0.0.0 libsaivs.so.0
	echo "$$DEB_CONTROL" > pkg/lucius-libsai/DEBIAN/control
	dpkg-deb --build pkg/lucius-libsai

lucius-libsai-bullseye:
	DOCKER_BUILDKIT=1 docker build . -f Dockerfile.saibuilder -t lemming-libsai:latest
	docker create --name libsai-temp lemming-libsai:latest
	docker cp libsai-temp:/build/pkg/lucius-libsai.deb .
	docker rm libsai-temp
