// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Авторизация пользователя и получение токена
	// (POST /auth/login)
	PostAuthLogin(c *gin.Context)
	// Регистрация пользователя
	// (POST /auth/register)
	PostAuthRegister(c *gin.Context)
	// Удаление баннеров по тегу или фиче
	// (DELETE /banner)
	DeleteBanner(c *gin.Context, params DeleteBannerParams)
	// Получение всех баннеров c фильтрацией по фиче и/или тегу
	// (GET /banner)
	GetBanner(c *gin.Context, params GetBannerParams)
	// Создание нового баннера
	// (POST /banner)
	PostBanner(c *gin.Context, params PostBannerParams)
	// Удаление баннера по идентификатору
	// (DELETE /banner/{id})
	DeleteBannerId(c *gin.Context, id int, params DeleteBannerIdParams)
	// Обновление содержимого баннера
	// (PATCH /banner/{id})
	PatchBannerId(c *gin.Context, id int, params PatchBannerIdParams)
	// Получение баннера для пользователя
	// (GET /user_banner)
	GetUserBanner(c *gin.Context, params GetUserBannerParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// PostAuthLogin operation middleware
func (siw *ServerInterfaceWrapper) PostAuthLogin(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostAuthLogin(c)
}

// PostAuthRegister operation middleware
func (siw *ServerInterfaceWrapper) PostAuthRegister(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostAuthRegister(c)
}

// DeleteBanner operation middleware
func (siw *ServerInterfaceWrapper) DeleteBanner(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params DeleteBannerParams

	// ------------- Optional query parameter "tag_id" -------------

	err = runtime.BindQueryParameter("form", true, false, "tag_id", c.Request.URL.Query(), &params.TagId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter tag_id: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "feature_id" -------------

	err = runtime.BindQueryParameter("form", true, false, "feature_id", c.Request.URL.Query(), &params.FeatureId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter feature_id: %w", err), http.StatusBadRequest)
		return
	}

	headers := c.Request.Header

	// ------------- Optional header parameter "token" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("token")]; found {
		var Token string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandler(c, fmt.Errorf("Expected one value for token, got %d", n), http.StatusBadRequest)
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "token", valueList[0], &Token, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter token: %w", err), http.StatusBadRequest)
			return
		}

		params.Token = &Token

	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.DeleteBanner(c, params)
}

// GetBanner operation middleware
func (siw *ServerInterfaceWrapper) GetBanner(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetBannerParams

	// ------------- Optional query parameter "feature_id" -------------

	err = runtime.BindQueryParameter("form", true, false, "feature_id", c.Request.URL.Query(), &params.FeatureId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter feature_id: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "tag_id" -------------

	err = runtime.BindQueryParameter("form", true, false, "tag_id", c.Request.URL.Query(), &params.TagId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter tag_id: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", c.Request.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter limit: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", c.Request.URL.Query(), &params.Offset)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter offset: %w", err), http.StatusBadRequest)
		return
	}

	headers := c.Request.Header

	// ------------- Optional header parameter "token" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("token")]; found {
		var Token string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandler(c, fmt.Errorf("Expected one value for token, got %d", n), http.StatusBadRequest)
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "token", valueList[0], &Token, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter token: %w", err), http.StatusBadRequest)
			return
		}

		params.Token = &Token

	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetBanner(c, params)
}

// PostBanner operation middleware
func (siw *ServerInterfaceWrapper) PostBanner(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostBannerParams

	headers := c.Request.Header

	// ------------- Optional header parameter "token" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("token")]; found {
		var Token string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandler(c, fmt.Errorf("Expected one value for token, got %d", n), http.StatusBadRequest)
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "token", valueList[0], &Token, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter token: %w", err), http.StatusBadRequest)
			return
		}

		params.Token = &Token

	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostBanner(c, params)
}

// DeleteBannerId operation middleware
func (siw *ServerInterfaceWrapper) DeleteBannerId(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params DeleteBannerIdParams

	headers := c.Request.Header

	// ------------- Optional header parameter "token" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("token")]; found {
		var Token string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandler(c, fmt.Errorf("Expected one value for token, got %d", n), http.StatusBadRequest)
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "token", valueList[0], &Token, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter token: %w", err), http.StatusBadRequest)
			return
		}

		params.Token = &Token

	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.DeleteBannerId(c, id, params)
}

// PatchBannerId operation middleware
func (siw *ServerInterfaceWrapper) PatchBannerId(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params PatchBannerIdParams

	headers := c.Request.Header

	// ------------- Optional header parameter "token" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("token")]; found {
		var Token string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandler(c, fmt.Errorf("Expected one value for token, got %d", n), http.StatusBadRequest)
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "token", valueList[0], &Token, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter token: %w", err), http.StatusBadRequest)
			return
		}

		params.Token = &Token

	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PatchBannerId(c, id, params)
}

// GetUserBanner operation middleware
func (siw *ServerInterfaceWrapper) GetUserBanner(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetUserBannerParams

	// ------------- Required query parameter "tag_id" -------------

	if paramValue := c.Query("tag_id"); paramValue != "" {

	} else {
		siw.ErrorHandler(c, fmt.Errorf("Query argument tag_id is required, but not found"), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "tag_id", c.Request.URL.Query(), &params.TagId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter tag_id: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Required query parameter "feature_id" -------------

	if paramValue := c.Query("feature_id"); paramValue != "" {

	} else {
		siw.ErrorHandler(c, fmt.Errorf("Query argument feature_id is required, but not found"), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "feature_id", c.Request.URL.Query(), &params.FeatureId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter feature_id: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "use_last_revision" -------------

	err = runtime.BindQueryParameter("form", true, false, "use_last_revision", c.Request.URL.Query(), &params.UseLastRevision)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter use_last_revision: %w", err), http.StatusBadRequest)
		return
	}

	headers := c.Request.Header

	// ------------- Optional header parameter "token" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("token")]; found {
		var Token string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandler(c, fmt.Errorf("Expected one value for token, got %d", n), http.StatusBadRequest)
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "token", valueList[0], &Token, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter token: %w", err), http.StatusBadRequest)
			return
		}

		params.Token = &Token

	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetUserBanner(c, params)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.POST(options.BaseURL+"/auth/login", wrapper.PostAuthLogin)
	router.POST(options.BaseURL+"/auth/register", wrapper.PostAuthRegister)
	router.DELETE(options.BaseURL+"/banner", wrapper.DeleteBanner)
	router.GET(options.BaseURL+"/banner", wrapper.GetBanner)
	router.POST(options.BaseURL+"/banner", wrapper.PostBanner)
	router.DELETE(options.BaseURL+"/banner/:id", wrapper.DeleteBannerId)
	router.PATCH(options.BaseURL+"/banner/:id", wrapper.PatchBannerId)
	router.GET(options.BaseURL+"/user_banner", wrapper.GetUserBanner)
}
