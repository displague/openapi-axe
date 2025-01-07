# openapi-axe

**openapi-axe** is a command-line tool to split a monolithic OpenAPI 3 specification into smaller, modular components. It organizes schemas, paths, and other reusable objects into a directory structure, making it easier to manage and maintain large OpenAPI specifications.

## Features

- Splits OpenAPI 3 specifications into structured directories.
- Supports extracting and organizing:
  - Components (e.g., schemas, examples, headers, request bodies, responses, and parameters).
  - Paths and their associated operations.
- Outputs clean YAML files for each extracted component.
- Updates references in the original specification to point to the split files.

## Installation

### Prerequisites

- Go 1.18 or newer.
- An OpenAPI 3 specification file to process.

### Build from Source

```bash
git clone https://github.com/displague/openapi-axe.git
cd openapi-axe
go build -o openapi-axe .
```

This will produce an `openapi-axe` binary in the current directory.

## Usage

### Command-Line Arguments

| Flag     | Description                            | Required |
|----------|----------------------------------------|----------|
| `-i`     | Path to the input OpenAPI specification file. | Yes      |
| `-d`     | Path to the target directory for output files. | Yes      |

### Example

Split a monolithic OpenAPI spec into a modular structure:

```bash
./openapi-axe -i ./openapi3.yaml -d ./output
```

### Output Structure

After running the tool, the target directory will be structured as follows:

```
output/
├── components/
│   ├── schemas/
│   │   ├── ExampleSchema1.yaml
│   │   └── ExampleSchema2.yaml
│   ├── examples/
│   ├── headers/
│   ├── requestBodies/
│   ├── responses/
│   └── parameters/
├── paths/
│   ├── /example1.yaml
│   └── /example2.yaml
└── openapi3.yaml
```

- Components (e.g., schemas) are saved in the corresponding subdirectory.
- Paths are saved as individual YAML files in the `paths` directory.
- References in the `openapi3.yaml` file are updated to point to the extracted files.

## Development

### Code Overview

- **Input Parsing**: The tool uses [`kin-openapi`](https://github.com/getkin/kin-openapi) to parse the OpenAPI specification.
- **File Handling**: Extracted components are written to files using Go's standard `os` and `io/fs` packages.
- **Directory Structure**: Components are organized into predefined subdirectories (e.g., `components/schemas`, `paths`).
- **Reference Updates**: Updates the OpenAPI document with references to the split files.

### Extending Functionality

To extend or customize the tool, you can:
- Add support for additional OpenAPI components (e.g., `securitySchemes`).
- Modify the directory structure or file-naming conventions.
- Implement additional validation or formatting.

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for enhancements or bug fixes.

## Acknowledgments

- Built using [kin-openapi](https://github.com/getkin/kin-openapi) for OpenAPI parsing and validation.
- Inspired by the need for modular management of complex OpenAPI specifications.
