#!/bin/bash

set -xe

cd ..
make deploy itest
make clean
make deploy2 itest2