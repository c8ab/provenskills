# provenskills

A CLI tool for packaging and managing skill definitions.

> **Work in progress** -- this project is under active development and not yet ready for general use.

## Requirements

- Go 1.23+

## Build

```sh
go build -o psk ./cmd/psk/
```

## Usage

```sh
# Validate a skill definition
psk validate ./path/to/skill-dir

# Build a skill package
psk build ./path/to/skill-dir --maintainer "Name <email>"

# List skills in the local store
psk list
```

## Development

```sh
# Run tests
go test ./...

# Run linter (requires golangci-lint v2)
golangci-lint run ./...
```

## License

See [LICENSE](LICENSE) for details.
