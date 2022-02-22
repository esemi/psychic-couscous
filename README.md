IMDB interactive graph
---

[wip-link](https://esemi.github.io/psychic-couscous/)


## Developers docs

### Build
```bash
$ sudo apt install golang-go
$ cd cmd/app/ && go build -o ../../bin && cd ../../
$ cp config.template.json bin/config.json
```

### Usage
```bash
$ bin/app help
Usage:
bin/app [command]

Available Commands:
completion  generate the autocompletion script for the specified shell
download    download imdb data
help        Help about any command
http-server Run http server
load        Loads data to neo4j
version     Application version

Flags:
-c, --config string   config file
-h, --help            help for bin/app
```
