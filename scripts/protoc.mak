.PHONY: proto openapi proto_deps clean_proto clean_openapi

PROTO_DIR := api/pkg/pb
OPENAPI_EXT := .swagger.json
PROTOC_INCLUDES := -I/usr/local/include -I$(PROTO_DIR) -I.
PROTO_SRC_FILES := $(shell find $(PROTO_DIR) -name \*.proto -not -path */google/* -not -path */grpc/* -not -path */protoc-gen-openapiv2/*)
OPENAPI_GEN_FILES := $(PROTO_SRC_FILES:.proto=$(OPENAPI_EXT))
PROTO_DEPENDENCY := $(PROTO_DIR)/google/api/annotations.proto \
	$(PROTO_DIR)/google/api/http.proto \
	$(PROTO_DIR)/grpc/validate/validate.proto \
	$(PROTO_DIR)/protoc-gen-openapiv2/options/annotations.proto \
	$(PROTO_DIR)/protoc-gen-openapiv2/options/openapiv2.proto

proto: $(PROTO_SRC_FILES:%.proto=%.pb.go)  ## generate go code for each proto files using protoc

openapi: $(OPENAPI_GEN_FILES)  ## generate openapi-spec for each proto files using protoc

proto_deps: $(PROTO_DEPENDENCY)  ## download proto dependencies

clean_proto:  ## delete generated protobuf code and downloaded files
	@echo "deleting proto generated and downloaded files..."
	@rm -f $(PROTO_SRC_FILES:%.proto=%*.pb.*go)
	@rm -f $(PROTO_DEPENDENCY)

clean_openapi:  ## delete generated openapi files
	@echo "deleting generated openapi spec files..."
	@rm -f $(OPENAPI_GEN_FILES)

%.pb.go: %.proto $(PROTO_DEPENDENCY)
	@echo "generating protobuf code for $<"
	@protoc $(PROTOC_INCLUDES) \
		--go_out=paths=source_relative:./$(PROTO_DIR) \
		--go-grpc_out=paths=source_relative:./$(PROTO_DIR) \
		--grpc-gateway_out=logtostderr=true,paths=source_relative:./$(PROTO_DIR) \
		--validate_out=lang=go,paths=source_relative:./$(PROTO_DIR) \
		$<
	@touch $@ # update timestamp of .pb.go even if the change is only in other files like validator and not in .pb.go

%$(OPENAPI_EXT): %.proto $(PROTO_DEPENDENCY)
	@echo "generating openapi definition for $<"
	@protoc $(PROTOC_INCLUDES) \
		--openapiv2_out $(PROTO_DIR) \
		--openapiv2_opt=logtostderr=true \
		$<

$(PROTO_DIR)/google/api/annotations.proto $(PROTO_DIR)/google/api/http.proto:
	@echo "downloading dependency file $@"
	@mkdir -p $(dir $@)
	@curl --fail -s -o $@ https://raw.githubusercontent.com/googleapis/googleapis/master/$(@:$(PROTO_DIR)/%=%)

$(PROTO_DIR)/grpc/validate/validate.proto:
	@echo "downloading dependency file $@"
	@mkdir -p $(dir $@)
	@curl --fail -s -o $@ https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/master/$(@:$(PROTO_DIR)/grpc/%=%)

$(PROTO_DIR)/protoc-gen-openapiv2/options/annotations.proto $(PROTO_DIR)/protoc-gen-openapiv2/options/openapiv2.proto:
	@echo "downloading dependency file $@"
	@mkdir -p $(dir $@)
	@curl --fail -s -o $@ https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/master/$(@:$(PROTO_DIR)/%=%)
