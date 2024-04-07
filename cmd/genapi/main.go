package genapi

//go:generate oapi-codegen -package api -generate types -o ../../internal/api/gen/types.gen.go ../../api/openapi.yaml
//go:generate oapi-codegen -package api -generate gin-server -o ../../internal/api/gen/server_base.gen.go ../../api/openapi.yaml
//go:generate oapi-codegen -package api -generate spec -o ../../internal/api/gen/server_spec.gen.go ../../api/openapi.yaml

// TODO: Сделать api и model генератор, который кладет все в папку gen в api и models
