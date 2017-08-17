# Myrmica Bibikoffi - Closes stale issues

[![Build Status](https://travis-ci.org/containous/bibikoffi.svg?branch=master)](https://travis-ci.org/containous/bibikoffi)
[![Docker Build Status](https://img.shields.io/docker/build/containous/bibikoffi.svg)](https://hub.docker.com/r/containous/bibikoffi/builds/)


```shell
Myrmica Bibikoffi: Closes stale issues.

Usage: bibikoffi [--flag=flag_argument] [-f[flag_argument]] ...     set flag_argument to flag(s)
   or: bibikoffi [--flag[=true|false| ]] [-f[true|false| ]] ...     set true/false to boolean flag(s)

Flags:
    --config-path Configuration file path.           (default "./bibikoffi.toml")
    --debug       Debug mode.                        (default "false")
    --dry-run     Dry run mode.                      (default "true")
-t, --token       GitHub Token. [required]           
-h, --help        Print Help (this message) and exit
```

## Description

Use a TOML configuration file. See `sample.toml`.

## Examples

```bash
bibikoffi -t xxxxxxxxxxxxxxx
```
