# Contributing

Issues and focused pull requests are welcome. Run the complete local checks
before submitting a pull request:

```bash
go mod tidy
go vet ./...
go test -race ./...
```

Public API changes must include tests and documentation. Never rewrite or move
an existing release tag; publish a new semantic version instead.
