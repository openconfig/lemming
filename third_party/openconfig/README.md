# Updating OpenConfig public Models in `lemming`

The OpenConfig modules are vendored into lemming to allow the implementation to
evolve at a pace that need not match HEAD of `openconfig/public`.

To update the submodule (and hence model version):

```bash
$ cd third_party/openconfig/public
$ git pull
$ git checkout <SHA or tag>`
```


