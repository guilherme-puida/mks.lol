# mks.lol

A privacy-oriented ephemeral link shortener.

---

Are you worried about malicious actors getting access to your chat longs and accessing private URLs that you've shared with friends and family?
Periodically deleting or encrypting these logs can guarantee that these urls won't be leaked, but in some cases you don't have this option.

By temporarily storing an url to another address, `mks.lol` guarantees that someone with access to your chat logs will not be able to access private information.

## Features

- Zero client-side javascript.
- Zero logging.
- Only depends on Go's standard library (excluding [styling with picocss](https://picocss.com)).
- In-memory storage only.
- Self-contained binary.

## Usage

There is a public instance running at [https://mks.lol/](https://mks.lol/), and pre-build binaries in the releases page.

### Building from source

```shell
$ git clone https://github.com/guilherme-puida/mks.lol.git
$ go build
$ ./mks.lol -h
Usage of ./mks.lol:
  -https
    	use https instead of http in rendered templates
  -port uint
    	port that will listen for all requests (default 8080)
  -url string
    	url used in rendered templates (default "mks.lol")
```