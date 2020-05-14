[![GitHub version](https://img.shields.io/github/v/release/jeffalyanak/check_freenas_api)](https://github.com/jeffalyanak/check_freenas_api/releases/latest)
[![License](https://img.shields.io/github/license/jeffalyanak/check_freenas_api)](https://github.com/jeffalyanak/check_freenas_api/blob/master/LICENSE)
[![Donate](https://img.shields.io/badge/donate--green)](https://jeff.alyanak.ca/donate)
[![Matrix](https://img.shields.io/badge/chat--green)](https://matrix.to/#/#check_freenas_api:social.rights.ninja)

# FreeNAS API Check Tool

Icinga/Nagios plugin that uses the FreeNAS API to check for alerts as well as pool health & usage/capacity.

Two types of check:
* alerts
* storage

Storage check has configurable percentage used `warning` and `critical` levels.

## Installation and requirements

If building from source:
* Golang 1.13.8

Otherwise, the binary can be used without any additional software.

## Usage

```bash
Usage of check_freenas_api:

Required:
  -check string
        Check to perform. Options are: {alerts,storage}
  -hostip string
        Host IP
  -username string
        Username
  -password string
        Password
  -skipverifytls
        Don't verify TLS certs

Optional:
  -crit int
        Storage used % for critical (default 90)
  -warn int
        Storage used % for warning (default 80
```

## License

FreeNAS API Check Tool is licensed under the terms of the GNU General Public License Version 3.
