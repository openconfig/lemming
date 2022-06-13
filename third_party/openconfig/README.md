# Updating OpenConfig public Models in `lemming`

The OpenConfig modules are vendored into lemming to allow the implementation to
evolve at a pace that need not match HEAD of `openconfig/public`.

To update this repo, pull the latest version from `openconfig/public` using
`git clone` and remove the `.git` directory. We avoid `git submodule` to avoid
the complexities of submodules, and the possibility of moving tags.


