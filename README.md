[![Build Status](https://www.travis-ci.com/roman-wb/websocket-mover.svg?branch=master)](https://www.travis-ci.com/roman-wb/websocket-mover)
![Go Report](https://goreportcard.com/badge/github.com/roman-wb/websocket-mover)
![Repository Top Language](https://img.shields.io/github/languages/top/roman-wb/websocket-mover)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/roman-wb/websocket-mover)
![Github Repository Size](https://img.shields.io/github/repo-size/roman-wb/websocket-mover)
![Lines of code](https://img.shields.io/tokei/lines/github/roman-wb/websocket-mover)
![License](https://img.shields.io/badge/license-MIT-green)
![GitHub last commit](https://img.shields.io/github/last-commit/roman-wb/websocket-mover)

# Move DOM object (Box) between browsers

Demo http://ec2-34-253-163-178.eu-west-1.compute.amazonaws.com/

## Feature

- Move DOM object (Box) between browsers
- Exclusive lock for move
- Toastr notify connect, disconnect, reconnect, etc

## Underhood

- WebSockets: gorilla/websocket
- Web framework: Echo
- JSON Logger: Echo

## Get Started

### Local server on localhost:8080

```bash
make server
```

### Mapping to the internet with ngrok

Note: Require installed ngrok https://ngrok.com

```bash
make ngrok
```

## Demo

## Inspired

https://github.com/gorilla/websocket/tree/master/examples/chat

## License

MIT
