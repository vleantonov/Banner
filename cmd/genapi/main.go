package genapi

//go:generate oapi-codegen -package api -generate types -o ../../internal/handler/http/v1/gen/types.gen.go ../../api/openapi.yaml
//go:generate oapi-codegen -package api -generate gin-server -o ../../internal/handler/http/v1/gen/server_base.gen.go ../../api/openapi.yaml
//go:generate oapi-codegen -package api -generate spec -o ../../internal/handler/http/v1/gen/server_spec.gen.go ../../api/openapi.yaml
