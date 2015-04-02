# findip
Small utility to find out a machine IP (as seen from the internet). 

## Installation

If you have a [Go(lang)](https://golang.org/) installation 

```
go get github.com/brunetto/findip
```

otherwise (if you run an Ubuntu 64bit) just download the binary.

## Usage

Launch with 

```bash
(./findip)&
```

It creates (and overwrite every 15 minutes) a `ip.txt` with the machine ip as seen from the internet.
