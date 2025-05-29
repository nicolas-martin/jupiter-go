#!/bin/bash

# Generate protobuf schemas from OpenAPI specs
openapi-generator generate \
	-i open-api/price.yaml \
	-g protobuf-schema \
	-o proto/price \
	--additional-properties=aggregateModelsName=models.proto \
	--additional-properties=numberedFieldNumberList=true \
	--additional-properties=startEnumsWithUnspecified=true \
	--additional-properties=supportMultipleResponses=true \
	--package-name=jupiter.price

# Fix the import path in the generated service file
# The generator incorrectly creates "models/models/proto.proto" instead of the actual filename
sed -i '' 's|models/models/proto\.proto|price/models/models_proto.proto|g' proto/price/services/default_service.proto

echo "Generated protobuf schemas and fixed import paths"

