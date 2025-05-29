openapi-generator generate \
  -i open-api/ultra.yaml \
  -g protobuf-schema \
  -o proto/ultra \
  --additional-properties=aggregateModelsName=combined_models.proto
