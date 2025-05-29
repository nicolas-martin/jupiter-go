# Jupiter Go Proto Client

This project generates Protocol Buffer client code for Jupiter REST endpoints using `buf` and OpenAPI specifications.

## Prerequisites

- [buf](https://buf.build/docs/installation) - Protocol buffer toolchain
- [openapi-generator](https://openapi-generator.tech/docs/installation) - OpenAPI code generator

### Why We Need These Fixes

The OpenAPI Generator (as of current versions) has several bugs when generating protobuf schemas that prevent successful compilation:

1. **❌ Incorrect Import Paths**
   - **Bug**: Generator creates `models/models/proto.proto` import paths
   - **Reality**: Actual file is `models/models_proto.proto`
   - **Fix**: `sed 's|models/models/proto\.proto|{service}/models/models_proto.proto|g'`

2. **❌ Unresolved AnyType References**
   - **Bug**: Generator leaves `AnyType` references in protobuf files
   - **Reality**: Should be `google.protobuf.Any`
   - **Fix**: `sed 's|AnyType|google.protobuf.Any|g'`

3. **❌ Double File Extensions**
   - **Bug**: Import mappings create `google/protobuf/any.proto.proto`
   - **Reality**: Should be `google/protobuf/any.proto`
   - **Fix**: `sed 's|google/protobuf/any\.proto\.proto|google/protobuf/any.proto|g'`

### Without These Fixes

Without the sed commands, you would get compilation errors like:
```bash
proto/service/services/default_service.proto:16:15:import "models/models/proto.proto": file does not exist
proto/service/models/models_proto.proto:25:15:field Service.extensions: unknown type AnyType
proto/service/models/models_proto.proto:15:15:import "google/protobuf/any.proto.proto": file does not exist
```

### Automatic Application

These fixes are automatically applied during `make generate-openapi` for every service, ensuring that:
- ✅ All protobuf files compile successfully
- ✅ Import paths are correct relative to the proto root
- ✅ Type mappings work properly
- ✅ No manual intervention is required

This approach allows us to use the OpenAPI Generator while working around its current limitations, providing a seamless development experience. 
