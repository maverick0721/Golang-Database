# Changelog

All notable changes to this project will be documented in this file.

The format is inspired by Keep a Changelog and this project follows Semantic Versioning.

## [Unreleased]

### Added
- Initial changelog structure.

### Changed
- N/A

### Fixed
- N/A

## [1.0.0] - 2026-03-22

### Added
- File-based JSON storage driver with collection/resource model.
- CRUD-style operations: Write, Read, ReadAll, Delete.
- Atomic write flow using temporary file rename.
- Per-collection mutex locking for safer concurrent writes/deletes.
- Demo program in main.
- Automated tests for core operations and validation paths.
- Enhanced README with architecture diagrams and usage documentation.

### Changed
- Go module version declaration made compatible with installed toolchain format.

### Fixed
- Multiple compile and runtime issues in driver methods and demo flow.
