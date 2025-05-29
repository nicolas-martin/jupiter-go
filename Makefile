.PHONY: help clean generate generate-proto generate-openapi install-deps

# Default target
help:
	@echo "Available targets:"
	@echo "  generate        - Generate all proto and OpenAPI client code"
	@echo "  generate-proto  - Generate proto client code using buf"
	@echo "  generate-openapi - Generate proto schemas from OpenAPI specs"
	@echo "  clean          - Clean generated files"
	@echo "  install-deps   - Install required dependencies"

# Generate all code
generate: generate-openapi generate-proto

# Generate proto client code using buf
generate-proto:
	@echo "Generating proto client code..."
	buf generate

# Generate proto schemas from OpenAPI specifications
generate-openapi:
	@echo "Generating proto schemas from OpenAPI specs..."
	@mkdir -p proto
	@for yaml_file in open-api/*.yaml; do \
		if [ -f "$$yaml_file" ]; then \
			service_name=$$(basename "$$yaml_file" .yaml); \
			echo "Generating protobuf schemas for $$service_name..."; \
			mkdir -p "proto/$$service_name"; \
			openapi-generator generate \
				-i "$$yaml_file" \
				-g protobuf-schema \
				-o "proto/$$service_name" \
				--additional-properties=aggregateModelsName=models.proto \
				--additional-properties=numberedFieldNumberList=true \
				--additional-properties=startEnumsWithUnspecified=true \
				--additional-properties=supportMultipleResponses=true \
				--package-name=jupiter.$$service_name \
				--type-mappings=AnyType=google.protobuf.Any \
				--import-mappings=google.protobuf.Any=google/protobuf/any.proto; \
			echo "Fixing import paths for $$service_name..."; \
			sed -i '' 's|models/models/proto\.proto|'$$service_name'/models/models_proto.proto|g' "proto/$$service_name/services"/*.proto 2>/dev/null || true; \
			sed -i '' 's|AnyType|google.protobuf.Any|g' "proto/$$service_name/models/models_proto.proto" 2>/dev/null || true; \
			sed -i '' 's|google/protobuf/any\.proto\.proto|google/protobuf/any.proto|g' "proto/$$service_name/models/models_proto.proto" 2>/dev/null || true; \
		fi; \
	done

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -rf gen/
	rm -rf proto/*

# Install required dependencies
install-deps:
	@echo "Installing dependencies..."
	@command -v buf >/dev/null 2>&1 || { echo "Installing buf..."; \
		curl -sSL "https://github.com/bufbuild/buf/releases/latest/download/buf-$$(uname -s)-$$(uname -m)" -o "/usr/local/bin/buf" && \
		chmod +x "/usr/local/bin/buf"; }
	@command -v openapi-generator >/dev/null 2>&1 || { echo "Please install openapi-generator: https://openapi-generator.tech/docs/installation"; exit 1; } 
