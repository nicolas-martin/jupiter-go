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
	@mkdir -p proto/ultra proto/swap proto/trigger proto/price proto/token proto/recurring
	openapi-generator generate \
		-i open-api/ultra.yaml \
		-g protobuf-schema \
		-o proto/ultra \
		--additional-properties=aggregateModelsName=combined_models.proto \
		--additional-properties=numberedFieldNumberList=true \
		--additional-properties=useProto3Optional=true \
		--type-mappings=AnyType=google.protobuf.Any \
		--import-mappings=google.protobuf.Any=google/protobuf/any.proto \
		--package-name=jupiter.ultra
	openapi-generator generate \
		-i open-api/swap.yaml \
		-g protobuf-schema \
		-o proto/swap \
		--additional-properties=aggregateModelsName=combined_models.proto \
		--additional-properties=numberedFieldNumberList=true \
		--additional-properties=useProto3Optional=true \
		--type-mappings=AnyType=google.protobuf.Any \
		--import-mappings=google.protobuf.Any=google/protobuf/any.proto \
		--package-name=jupiter.swap
	openapi-generator generate \
		-i open-api/trigger.yaml \
		-g protobuf-schema \
		-o proto/trigger \
		--additional-properties=aggregateModelsName=combined_models.proto \
		--additional-properties=numberedFieldNumberList=true \
		--additional-properties=useProto3Optional=true \
		--type-mappings=AnyType=google.protobuf.Any \
		--import-mappings=google.protobuf.Any=google/protobuf/any.proto \
		--package-name=jupiter.trigger
	openapi-generator generate \
		-i open-api/price.yaml \
		-g protobuf-schema \
		-o proto/price \
		--additional-properties=aggregateModelsName=combined_models.proto \
		--additional-properties=numberedFieldNumberList=true \
		--additional-properties=useProto3Optional=true \
		--type-mappings=AnyType=google.protobuf.Any \
		--import-mappings=google.protobuf.Any=google/protobuf/any.proto \
		--package-name=jupiter.price
	openapi-generator generate \
		-i open-api/token.yaml \
		-g protobuf-schema \
		-o proto/token \
		--additional-properties=aggregateModelsName=combined_models.proto \
		--additional-properties=numberedFieldNumberList=true \
		--additional-properties=useProto3Optional=true \
		--type-mappings=AnyType=google.protobuf.Any \
		--import-mappings=google.protobuf.Any=google/protobuf/any.proto \
		--package-name=jupiter.token
	openapi-generator generate \
		-i open-api/recurring.yaml \
		-g protobuf-schema \
		-o proto/recurring \
		--additional-properties=aggregateModelsName=combined_models.proto \
		--additional-properties=numberedFieldNumberList=true \
		--additional-properties=useProto3Optional=true \
		--type-mappings=AnyType=google.protobuf.Any \
		--import-mappings=google.protobuf.Any=google/protobuf/any.proto \
		--package-name=jupiter.recurring

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