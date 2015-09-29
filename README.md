# Atoll

A simple JSON friendly system monitoring agent.

## Build

```bash
make
```

## Run

```bash
./bin/atoll
```

## Test

```bash
make test
make test.verbose
```

## Example: Docker Web app

```bash
cd examples/docker/web-app-small
make
make run
```

## TODO

* Proper binary packages (CentOS/Redhat, Ubuntu/Debian, etc.)
* More comprehensive Docker based simulations (ETL, Hadoop ecosystem, etc.)
* More elegant Docker integration (run on Docker boot?)
* Detached service specific modules: `apt-get install atoll-elasticsearch`
