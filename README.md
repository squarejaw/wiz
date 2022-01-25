# wiz
A CLI tool (written in Go) to control Phillips WiZ light bulbs

Inspired by [pywizlight](https://github.com/sbidy/pywizlight).

## Installation
```
go install github.com/squarejaw/wiz@main
```

Or download a binary [release](https://github.com/squarejaw/wiz/releases).

## Usage
```sh
$ wiz list

IP              MAC
192.168.1.12    6ab731ba1bd5
192.168.1.34    86082d45c947

$ wiz on 192.168.1.12 -k 2700 -d 50
{"success":true}

$ wiz state 192.168.1.12
{"mac":"6ab731ba1bd5","rssi":-37,"src":"","state":true,"sceneId":0,"temp":2700,"dimming":50}

$ wiz off 192.168.1.12
{"success":true}
```
