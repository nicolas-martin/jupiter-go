.PHONY: help generate generate-openapi generate-proto clean install-deps run-client

help: ## Show available targets
	@echo "Available targets:"
	@echo "  help              Show this help message"
	@echo "  generate          Generate all proto and OpenAPI client code"
	@echo "  generate-openapi  Generate proto schemas from OpenAPI specs"
	@echo "  generate-proto    Generate client code using buf"
	@echo "  run-client        Run the Jupiter client demo"
	@echo "  clean             Clean generated files"
	@echo "  install-deps      Install required dependencies"

generate: generate-openapi generate-proto ## Generate all proto and OpenAPI client code

generate-openapi: ## Generate proto schemas from OpenAPI specs
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

generate-proto: ## Generate client code using buf
	@echo "Generating proto client code..."
	buf generate

run-client: ## Run the Jupiter client demo
	@echo "Running Jupiter client demo..."
	go run cmd/main.go

clean: ## Clean generated files
	@echo "Cleaning generated files..."
	rm -rf gen/
	rm -rf proto/*

install-deps: ## Install required dependencies
	@echo "Installing dependencies..."
	@echo "Installing buf..."
	@if ! command -v buf >/dev/null 2>&1; then \
		echo "Please install buf from https://buf.build/docs/installation"; \
		exit 1; \
	else \
		echo "✓ buf is already installed"; \
	fi
	@echo "Installing openapi-generator..."
	@if ! command -v openapi-generator >/dev/null 2>&1; then \
		echo "Please install openapi-generator from https://openapi-generator.tech/docs/installation"; \
		exit 1; \
	else \
		echo "✓ openapi-generator is already installed"; \
	fi
	@echo "Installing Go dependencies..."
	go mod tidy
	@echo "✓ All dependencies are ready" 
