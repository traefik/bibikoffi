# Myrmica Bibikoffi - Closes stale issues

[![GitHub release](https://img.shields.io/github/release/containous/bibikoffi.svg)](https://github.com/containous/bibikoffi/releases/latest)
[![Build Status](https://travis-ci.com/containous/bibikoffi.svg?branch=master)](https://travis-ci.com/containous/bibikoffi)
[![Docker Build Status](https://img.shields.io/docker/build/containous/bibikoffi.svg)](https://hub.docker.com/r/containous/bibikoffi/builds/)

## Description

Use a TOML configuration file. See [sample.toml](/sample.toml).

```yaml
Myrmica Bibikoffi: Closes stale issues.

Usage: bibikoffi [--flag=flag_argument] [-f[flag_argument]] ...     set flag_argument to flag(s)
   or: bibikoffi [--flag[=true|false| ]] [-f[true|false| ]] ...     set true/false to boolean flag(s)

Flags:
    --config-path Configuration file path.           (default "./bibikoffi.toml")
    --debug       Debug mode.                        (default "false")
    --dry-run     Dry run mode.                      (default "true")
    --port        Server port.                       (default "80")
    --server      Server mode.                       (default "false")
-t, --token       GitHub Token. [required]           
-h, --help        Print Help (this message) and exit
```

## Examples

```bash
bibikoffi -t xxxxxxxxxxxxxxx
```

## What does Myrmica Bibikoffi mean?

![Myrmica Bibikoffi](http://www.antwiki.org/wiki/images/2/28/Myrmica_bibikoffi_H_casent0900283.jpg)
