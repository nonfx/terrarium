PROTO_SRC_EXT = .proto
PROTO_GEN_EXT = .pb.go
PROTO_GEN_GW_EXT = .pb.gw.go
PROTO_GEN_GRPC_EXT = _grpc.pb.go
PROTO_GEN_VALIDATOR_EXT = .pb.validate.go
PROTO_SRC_FILES := $(shell find api/pkg -name \*${PROTO_SRC_EXT} -not -path */google/* -not -path */grpc/* -not -path */protoc-gen-openapiv2/*)
PROTO_SRC_FILES_ALL := $(shell find api/pkg -name \*${PROTO_SRC_EXT})
PROTO_GEN_FILES := $(PROTO_SRC_FILES:$(PROTO_SRC_EXT)=$(PROTO_GEN_EXT))
PROTO_GEN_FILES_ALL := $(PROTO_SRC_FILES_ALL:$(PROTO_SRC_EXT)=$(PROTO_GEN_EXT))
PROTO_GEN_GW_FILES := $(PROTO_SRC_FILES_ALL:$(PROTO_SRC_EXT)=$(PROTO_GEN_GW_EXT))
PROTO_GEN_VALIDATOR_FILES := $(PROTO_SRC_FILES_ALL:$(PROTO_SRC_EXT)=$(PROTO_GEN_VALIDATOR_EXT))
OPENAPI_GEN_FILES := $(PROTO_SRC_FILES:$(PROTO_SRC_EXT)=$(OPENAPI_EXT))
OPENAPI_GEN_FILES_ALL := $(PROTO_SRC_FILES_ALL:$(PROTO_SRC_EXT)=$(OPENAPI_EXT))
PROTO_GEN_GRPC_FILES := $(PROTO_SRC_FILES_ALL:$(PROTO_SRC_EXT)=$(PROTO_GEN_GRPC_EXT))

GATEWAY_DEPENDENCIES := api/pkg/pb/google/api/annotations.proto api/pkg/pb/google/api/http.proto
GATEWAY_DEPENDENCY := $(shell $(foreach f,$(GATEWAY_DEPENDENCIES),test -e "$f" || echo "$f";))

GRPC_VALIDATOR_FILE := api/pkg/pb/grpc/validator/validator.proto
GRPC_VALIDATOR_DEPENDENCY := $(shell $(foreach f,$(GRPC_VALIDATOR_FILE),test -e "$f" || echo "$f";))

OPENAPI_DEPENDENCIES := api/pkg/pb/protoc-gen-openapiv2/options/openapiv2.proto api/pkg/pb/protoc-gen-openapiv2/options/annotations.proto
OPENAPI_DEPENDENCY := $(shell $(foreach f,$(OPENAPI_DEPENDENCIES),test -e "$f" || echo "$f";))

define proto_gen_rule
$(1)$(PROTO_GEN_EXT): $(1)$(PROTO_SRC_EXT) | $(GATEWAY_DEPENDENCY) $(GRPC_VALIDATOR_DEPENDENCY) $(PROTOC_GEN_GO) $(PROTOC_GEN_GRPC_GATEWAY) $(PROTOC_GEN_GPRC_GO) $(OPENAPI_DEPENDENCY)
	@echo "\033[32m-- Generating protobuf code for $$< \033[0m"
	protoc --validate_out=lang=go,paths=source_relative:./api/pkg/pb -I/usr/local/include -I./api/pkg/pb -I. --plugin=protoc-gen-grpc=grpc_go_plugin --grpc-gateway_out=logtostderr=true,paths=source_relative:./api/pkg/pb --go-grpc_out=paths=source_relative:./api/pkg/pb --go_out=paths=source_relative:./api/pkg/pb $$<
endef

define openapi_gen_rule
$(1)$(OPENAPI_EXT): $(1)$(PROTO_SRC_EXT) | $(PROTOC_GEN_OPENAPIV2) $(SWAGGER_CODEGEN) $(OPENAPI_DEPENDENCY) $(GATEWAY_DEPENDENCY) $(GRPC_VALIDATOR_DEPENDENCY)
	@echo "\033[32m-- Generating OpenAPI definition for $$< \033[0m"
	protoc --validate_out=lang=go,paths=source_relative:./api/pkg/pb -I/usr/local/include -I./api/pkg/pb -I. --openapiv2_out api/pkg/pb --openapiv2_opt logtostderr=true $$<
endef

$(foreach proto_file,$(PROTO_SRC_FILES),$(eval $(call proto_gen_rule,$(proto_file:.proto=))))

$(foreach proto_file,$(PROTO_SRC_FILES),$(eval $(call openapi_gen_rule,$(proto_file:.proto=))))

.PHONY: proto
proto: $(PROTO_GEN_FILES)  ## Generate Go code from protocol buffer definitions

.PHONY: openapi
openapi: $(OPENAPI_GEN_FILES)  ## Generate OpenAPI JSON files from protocol buffer definition

$(GATEWAY_DEPENDENCY):
	mkdir -p api/pkg/pb/google/api
	@curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto > api/pkg/pb/google/api/annotations.proto
	@curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto > api/pkg/pb/google/api/http.proto

$(GRPC_VALIDATOR_DEPENDENCY):
	mkdir -p api/pkg/pb/grpc/validator
	@curl https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/master/validate/validate.proto > api/pkg/pb/grpc/validator/validator.proto

$(OPENAPI_DEPENDENCY):
	mkdir -p api/pkg/pb/protoc-gen-openapiv2/options
	@curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/master/protoc-gen-openapiv2/options/annotations.proto > api/pkg/pb/protoc-gen-openapiv2/options/annotations.proto
	@curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/master/protoc-gen-openapiv2/options/openapiv2.proto > api/pkg/pb/protoc-gen-openapiv2/options/openapiv2.proto

clean_proto:  ## Delete the generated protobuf files
	rm -f $(PROTO_GEN_FILES_ALL)
	rm -f $(PROTO_GEN_GW_FILES)
	rm -f $(PROTO_GEN_VALIDATOR_FILES)
	rm -f $(PROTO_GEN_GRPC_FILES)
	rm -f $(GATEWAY_DEPENDENCIES)
	rm -f $(GRPC_VALIDATOR_FILE)
	rm -f $(OPENAPI_DEPENDENCIES)

clean_openapi:  ## Delete the generated OpenAPI files
	rm -f $(OPENAPI_GEN_FILES_ALL)
	rm -f $(OPENAPI_DEPENDENCIES)
