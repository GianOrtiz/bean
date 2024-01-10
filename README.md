# Bean

Bean is an abstract implementation of a double entry journal account system for experimentation. It only provides a backend that can be used by any system.

## Requirements

* Go +v1.21
* [Migrate command line tool](https://github.com/golang-migrate/migrate/tree/master)
* SQLite3
* [mockgen](https://github.com/uber-go/mock)

## Running

1. Build the release with the command `make build`.
2. Run the server with the command `./main`.
