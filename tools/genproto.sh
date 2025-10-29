#!/bin/bash
set -xe

gendir=$(bazel info bazel-genfiles)

proto_libs=$(bazel query 'kind(go_proto_library, //...)')
for fulltarget in $proto_libs; do
    bazel build "$fulltarget"
    importpath=$(bazel query "$fulltarget" --output=build | sed -n 's/.*importpath = "\(.*\)".*/\1/p')
    dir=$(sed -E 's/\/\/(.*):.*/\1/g' <<<"$fulltarget")
    target=$(sed -E 's/.*:(.*)_go_proto/\1/g' <<<"$fulltarget")
    parentdir=$(dirname "$dir")
    # Files in bazel-bin are not writable by default and need to be for new files.
    chmod u+w -R "$gendir"/"$dir"/"$target"_go_proto_/"$importpath"
    cp -fr "$gendir"/"$dir"/"$target"_go_proto_/"$importpath" "$parentdir"
done
