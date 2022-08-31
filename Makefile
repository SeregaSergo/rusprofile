NAME 		:= server
PROTO_NAME	:= api.proto
API_PATH	:= ./proto
PROTO_FILE 	:= $(API_PATH)/$(PROTO_NAME)
.PHONY: all bin api dep

all: bin api dep $(NAME) ## Build the binary file for server

$(NAME):
	go build -o $(NAME) ./main.go

api: $(PROTO_FILE)
	protoc -I ./proto --go_out=. \
					--go-grpc_out=. \
					--grpc-gateway_out=. \
					--grpc-gateway_opt logtostderr=true \
					--grpc-gateway_opt generate_unbound_methods=true \
					--openapiv2_out=swagger-ui $(PROTO_FILE)

dep: ## Get the dependencies
	go get -v -d ./...

bin:
	go install "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway" \
				"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2" \
				"google.golang.org/grpc/cmd/protoc-gen-go-grpc" \
				"google.golang.org/protobuf/cmd/protoc-gen-go"

clean: ## Remove previous builds
	@rm $(NAME) ./internal/api/*

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
