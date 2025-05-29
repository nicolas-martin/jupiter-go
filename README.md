# Jupiter Go Proto Client

This project generates Protocol Buffer client code for Jupiter REST endpoints using `buf` and OpenAPI specifications.

## Prerequisites

- [buf](https://buf.build/docs/installation) - Protocol buffer toolchain
- [openapi-generator](https://openapi-generator.tech/docs/installation) - OpenAPI code generator

## Quick Start

1. **Install dependencies:**
   ```bash
   make install-deps
   ```

2. **Generate all client code:**
   ```bash
   make generate
   ```

## Available Commands

| Command | Description |
|---------|-------------|
| `make help` | Show available targets |
| `make generate` | Generate all proto and OpenAPI client code |
| `make generate-openapi` | Generate proto schemas from OpenAPI specs |
| `make generate-proto` | Generate client code using buf |
| `make clean` | Clean generated files |
| `make install-deps` | Install required dependencies |

## Project Structure

```
├── open-api/           # OpenAPI specifications
│   ├── ultra.yaml
│   ├── swap.yaml
│   ├── trigger.yaml
│   ├── price.yaml
│   ├── token.yaml
│   └── recurring.yaml
├── proto/              # Generated proto schemas
├── gen/                # Generated client code
│   ├── go/            # Go client code
│   ├── js/            # JavaScript client code
│   └── python/        # Python client code
├── buf.yaml           # Buf configuration
├── buf.gen.yaml       # Buf generation configuration
└── Makefile           # Build automation
```

## Generated Output

- **Proto schemas**: Generated in `proto/` directories from OpenAPI specs
- **Client code**: Generated in `gen/` with support for Go, JavaScript, and Python
- **Combined models**: Each service gets a `combined_models.proto` file with all models

## Configuration

- `buf.yaml`: Buf workspace configuration
- `buf.gen.yaml`: Code generation plugins and output settings
- OpenAPI specs are converted to proto schemas with numbered fields and proto3 optional support 