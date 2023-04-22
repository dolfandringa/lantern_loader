# Lantern Loader
[![Unittesting](https://github.com/dolfandringa/lantern_loader/actions/workflows/go.yml/badge.svg)](https://github.com/dolfandringa/lantern_loader/actions/workflows/go.yml)


## Compiling
`go build -o bin`

## Testing
`go test -v ./...`

## Running

`lantern_loader url1 url2...`

This will download a single file from multiple sources as specified in the urls. 

example:

```
./bin/lantern_loader https://mirror.pit.teraswitch.com/ubuntu-releases/23.04/ubuntu-23.04-desktop-amd64.iso https://southfront.mm.fcix.net/ubuntu-releases/23.04/ubuntu-23.04-desktop-amd64.iso  https://forksystems.mm.fcix.net/ubuntu-releases/23.04/ubuntu-23.04-desktop-amd64.iso https://mirror.math.princeton.edu/pub/ubuntu-iso/23.04/ubuntu-23.04-desktop-amd64.iso

File size:  4932407296
Dowloading  471 chunks
Starting worker for url https://mirror.pit.teraswitch.com/ubuntu-releases/23.04/ubuntu-23.04-desktop-amd64.iso
Starting worker for url https://southfront.mm.fcix.net/ubuntu-releases/23.04/ubuntu-23.04-desktop-amd64.iso
Starting worker for url https://forksystems.mm.fcix.net/ubuntu-releases/23.04/ubuntu-23.04-desktop-amd64.iso
Starting worker for url https://mirror.math.princeton.edu/pub/ubuntu-iso/23.04/ubuntu-23.04-desktop-amd64.iso
   9% |█████████                                                                                                    | (430 MB/4.6 GB, 24 MB/s) [22s:2m59s]
```
