---
name: installing
description: Installing gotestproxy on a system, so it is available in the PATH.
---

# Installing on a system

To install `gotestproxy` binary, simply run `make install`,
setting the `PREFIX` env variable to a directory that is in your `PATH` env variable.

```sh
PREFIX=<target_directory> make install
```

Uninstalling is also very simple:

```sh
PREFIX=<previously_set_prefix> make uninstall
```

Alternatively, you can simply remove the `gotestproxy` binary — it's self-contained,
so it's the only file to remove.
