#!/bin/sh

# `$*` expands the `args` supplied in an `array` individually
# or splits `args` in a string separated by whitespace.
cd /tmp/ || exit 1
sh -c "/app/webdl $*"
