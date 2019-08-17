// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/pseudo-su/oapi-codegen DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/pseudo-su/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// (GET /pets)
	FindPets(ctx echo.Context, params FindPetsParams) error
	// (POST /pets)
	AddPet(ctx echo.Context) error
	// (DELETE /pets/{id})
	DeletePet(ctx echo.Context, id int64) error
	// (GET /pets/{id})
	FindPetById(ctx echo.Context, id int64) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// FindPets converts echo context to params.
func (w *ServerInterfaceWrapper) FindPets(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params FindPetsParams
	// ------------- Optional query parameter "tags" -------------
	if paramValue := ctx.QueryParam("tags"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "tags", ctx.QueryParams(), &params.Tags)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter tags: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------
	if paramValue := ctx.QueryParam("limit"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.FindPets(ctx, params)
	return err
}

// AddPet converts echo context to params.
func (w *ServerInterfaceWrapper) AddPet(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AddPet(ctx)
	return err
}

// DeletePet converts echo context to params.
func (w *ServerInterfaceWrapper) DeletePet(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameter("simple", false, "id", ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeletePet(ctx, id)
	return err
}

// FindPetById converts echo context to params.
func (w *ServerInterfaceWrapper) FindPetById(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameter("simple", false, "id", ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.FindPetById(ctx, id)
	return err
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router runtime.EchoRouter, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET("/pets", wrapper.FindPets)
	router.POST("/pets", wrapper.AddPet)
	router.DELETE("/pets/:id", wrapper.DeletePet)
	router.GET("/pets/:id", wrapper.FindPetById)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}
