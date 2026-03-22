# Golang-Database

<p align="center">
	<img alt="Golang-Database Banner" src="https://capsule-render.vercel.app/api?type=waving&height=220&color=0:0B132B,50:1C2541,100:3A506B&text=Golang-Database&fontColor=EAF4F4&fontSize=48&fontAlignY=40&desc=Lightweight%20JSON%20Database%20Engine%20in%20Go&descAlignY=62" />
</p>

<p align="center">
	<a href="https://github.com/maverick0721/Golang-Database">
		<img alt="Repository" src="https://img.shields.io/badge/GitHub-Repository-0B132B?style=for-the-badge&logo=github&logoColor=white">
	</a>
	<a href="https://github.com/maverick0721/Golang-Database/stargazers">
		<img alt="GitHub Stars" src="https://img.shields.io/github/stars/maverick0721/Golang-Database?style=for-the-badge&color=3A506B&logo=github&logoColor=white">
	</a>
	<a href="https://github.com/maverick0721/Golang-Database/issues">
		<img alt="GitHub Issues" src="https://img.shields.io/github/issues/maverick0721/Golang-Database?style=for-the-badge&color=5BC0BE">
	</a>
	<img alt="Go Version" src="https://img.shields.io/badge/Go-1.18%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white">
	<img alt="Storage" src="https://img.shields.io/badge/Storage-JSON%20Files-2F855A?style=for-the-badge">
	<img alt="Concurrency" src="https://img.shields.io/badge/Concurrency-Mutex%20Protected-1C7C7D?style=for-the-badge">
	<img alt="Tests" src="https://img.shields.io/badge/Tests-go%20test%20.%2F...-22863A?style=for-the-badge">
</p>

<p align="center">
	Lightweight JSON file database engine in Go for local-first apps, tools, and prototypes.
</p>

## Table of Contents

- [Overview](#overview)
- [Project Status](#project-status)
- [Highlights](#highlights)
- [Quick Start](#quick-start)
- [CLI Commands](#cli-commands)
- [Architecture](#architecture)
- [Write Flow (Atomic)](#write-flow-atomic)
- [Data Layout](#data-layout)
- [Demo Output](#demo-output)
- [API Snapshot](#api-snapshot)
- [Usage Example](#usage-example)
- [Testing](#testing)
- [Changelog](#changelog)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [Notes](#notes)

## Overview

This project stores records as JSON files on disk, uses temp-file rename for safer writes, and protects writes/deletes with per-collection mutexes.

It is intentionally simple and easy to inspect, making it a good fit for:

- local tools and CLI prototypes
- educational projects
- lightweight persistence where running a full DB server is unnecessary

## Project Status

- Current scope: single-process, file-based JSON storage
- Verified locally with tests and demo run
- API surface is stable for core operations (`Write`, `Read`, `ReadAll`, `Delete`)

## Highlights

- Zero external database setup
- Atomic file writes (`.tmp` -> `.json`)
- Collection-scoped locking for concurrency safety
- CRUD-style operations (`Write`, `Read`, `ReadAll`, `Delete`)
- Simple API for local tools, demos, and small apps

## Quick Start

```bash
go mod tidy
go test ./...
go run .
```

## CLI Commands

```bash
# Run demo program
go run .

# Build
go build ./...

# Run tests
go test ./...
```

## Architecture

```mermaid
flowchart TD
	A[Application Code] --> B[Driver]
	B --> C{Operation}
	C -->|Write| D[Lock collection mutex]
	D --> E[Marshal JSON]
	E --> F[Write *.tmp file]
	F --> G[Rename to *.json]
	C -->|Read| H[Read single .json file]
	C -->|ReadAll| I[List collection dir and read files]
	C -->|Delete| J[Lock collection mutex and remove file]
	G --> K[(Filesystem)]
	H --> K
	I --> K
	J --> K
```

## Write Flow (Atomic)

```mermaid
sequenceDiagram
	participant App as App
	participant DB as Driver
	participant FS as Filesystem

	App->>DB: Write(collection, resource, payload)
	DB->>DB: Validate inputs
	DB->>DB: Lock collection mutex
	DB->>FS: mkdir -p data/collection
	DB->>FS: write resource.json.tmp
	DB->>FS: rename tmp -> resource.json
	DB-->>App: success
```

## Data Layout

After running the demo, records are stored like:

```text
data/
  users/
	Aman.json
	Manav.json
	Priyanshu.json
	Shailendra.json
	Siddharth.json
	Yash.json
```

## Demo Output

Sample output from `go run .`:

```text
[{
	"Name": "Aman",
	"Age": 21,
	"Contact": "0987651111",
	"Company": "Apple",
	"Address": {
		"City": "Gwalior",
		"State": "Madhya Pradesh",
		"Country": "India",
		"Pincode": 474011
	}
}
 ...
]
[{Aman 21 0987651111 Apple {Gwalior Madhya Pradesh India 474011}} ...]
```

## API Snapshot

Core methods:

```go
Write(collection, resource string, v interface{}) error
Read(collection, resource string, v interface{}) error
ReadAll(collection string) ([]string, error)
Delete(collection, resource string) error
```

## Usage Example

This project is currently a `main` package, so usage is demonstrated in `main.go`.

```go
db, err := New("./data", nil)
if err != nil {
	panic(err)
}

u := User{Name: "Alice", Age: "21", Contact: "9999999999", Company: "ExampleCo"}

if err := db.Write("users", u.Name, u); err != nil {
	panic(err)
}

var out User
if err := db.Read("users", "Alice", &out); err != nil {
	panic(err)
}

records, err := db.ReadAll("users")
if err != nil {
	panic(err)
}

_ = records
```

## Testing

Automated tests currently cover:

- Write/Read round-trip
- ReadAll returns all records
- Delete removes records
- Validation for empty collection/resource inputs

Run tests:

```bash
go test ./...
```

## Changelog

Release notes are tracked in [CHANGELOG.md](CHANGELOG.md).

## Roadmap

- Add benchmarks for write/read/delete throughput
- Add package split (`driver` package) for easier import into other projects
- Optional indexing layer for faster lookups in larger collections
- Better error taxonomy with sentinel errors

## Contributing

Contributions are welcome.

1. Fork the repository
2. Create a feature branch
3. Add or update tests
4. Run `go test ./...`
5. Open a pull request with a clear summary

## Notes

- Recommended for lightweight/local persistence use cases.
- Not intended as a replacement for full production RDBMS/NoSQL systems.
