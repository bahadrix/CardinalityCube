![Logo](https://user-images.githubusercontent.com/609775/79579888-dc8cf680-80d0-11ea-838e-603cebd8c00c.png)

[![Build Status](https://travis-ci.org/bahadrix/cardinalitycube.svg?branch=master)](https://travis-ci.org/bahadrix/cardinalitycube) [![Go Report Card](https://goreportcard.com/badge/github.com/bahadrix/cardinalitycube)](https://goreportcard.com/report/github.com/bahadrix/cardinalitycube) [![codecov](https://codecov.io/gh/bahadrix/cardinalitycube/branch/master/graph/badge.svg)](https://codecov.io/gh/bahadrix/cardinalitycube) [![Reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/bahadrix/cardinalitycube?tab=subdirectories) 


# Cardinality Cube 
Fast and Accurate Approximate Cardinality Estimator Data Structure and Server

## Dependencies

Following packages needed by GOCMZQ

- pkg-config
- libczmq-dev
- libsodium-dev

You can install them in Debian based system like:
```bash
apt-get install -y pkg-config libczmq4 libczmq-dev libsodium-dev
```


# Possible Usages
Thanks to modular design of CC ecosystem you can either;
- Use the cube as a data structure in your code *(see readme file under cube folder)*
- Embed the server in your code *(see `server/service/cmd/start.go` for usage sample)*
- Use the full blown server as standalone *(you are here)*

these options also available for client usage:
