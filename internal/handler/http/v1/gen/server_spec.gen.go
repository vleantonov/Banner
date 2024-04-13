// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"KVJKJAAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
