#!/bin/bash

gendir=$(bazel info bazel-genfiles)
importpath="github.com/openconfig/lemming"

proto_libs=$(bazel query 'kind(go_proto_library, //...)')
for fulltarget in $proto_libs; do
    bazel build $fulltarget
    dir=$(sed -E 's/\/\/(.*):.*/\1/g' <<<"$fulltarget")
    target=$(sed -E 's/.*:(.*)_go_proto/\1/g' <<<"$fulltarget")
    cp -r "$gendir/$dir/$target"_go_proto_/"$importpath/$dir" proto
done