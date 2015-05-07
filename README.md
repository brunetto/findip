# fip
Small utility to find out a machine IP (as seen from the internet). 

## Installation

If you have a [Go(lang)](https://golang.org/) installation 

```
go get github.com/brunetto/fip
```

otherwise (if you run an Ubuntu 64bit) just download the binary.

## Usage

Use like:

```
fip once
```

to run it one time and get back the ip on the STDOUT 
(useful to start a file server for example)

Launch with 

```bash
(./fip)&
```

to have it run indefinitely and updating the ip.txt file
(useful to have an update address to connect to).

