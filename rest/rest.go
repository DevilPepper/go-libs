package rest

import (
	"fmt"

	"github.com/DevilPepper/go-libs/environment"
	"github.com/DevilPepper/go-libs/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/goccy/go-yaml"
)

func NewServer(enableSwagger bool, port int, withInfo func(e *fuego.Engine), serverOptions ...func(*fuego.Server)) *fuego.Server {
	host := ""
	if environment.IS_DEV {
		host = environment.GetBaseUrl()
		if host == "" {
			host = "localhost"
		}
	}

	opts := append(
		[]func(*fuego.Server){
			fuego.WithEngineOptions(
				fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
					Disabled:         !enableSwagger,
					DisableSwaggerUI: !enableSwagger,
					DisableLocalSave: true,
					DisableMessages:  false,
					PrettyFormatJSON: false,
					SwaggerURL:       "/api/docs",
					SpecURL:          "/api/docs/openapi.json",
					UIHandler:        fuego.DefaultOpenAPIHandler,
					MiddlewareConfig: fuego.MiddlewareConfig{
						DisableMiddlewareSection: true,
					},
				}),
				withInfo,
			),
			fuego.WithRouteOptions(
				option.DefaultStatusCode(200),
			),
			fuego.WithLoggingMiddleware(fuego.LoggingConfig{
				DisableRequest:  environment.IS_DEV,
				DisableResponse: environment.IS_DEV,
			}),
			fuego.WithAddr(fmt.Sprintf("%s:%d", host, port)),
		},
		serverOptions...,
	)

	s := fuego.NewServer(
		opts...,
	)
	if environment.IS_DEV {
		fuego.Use(s, middleware.RequestLogs)
	}
	return s
}

func WithInfo(title string, description string, version string) func(e *fuego.Engine) {
	return func(e *fuego.Engine) {
		info := e.OpenAPI.Description().Info
		info.Title = title
		info.Description = description
		info.Version = version
	}
}

func fixOperation(op *openapi3.Operation) {
	if op != nil {
		ok := op.Responses.Map()["204"]
		if ok != nil {
			ok.Value.Content = make(openapi3.Content)
		}
	}
}

func FixSpec(s *fuego.Server) {
	spec := s.OpenAPI.Description()
	delete(spec.Components.Schemas, "unknown-interface")
	for _, path := range spec.Paths.Map() {
		fixOperation(path.Get)
		fixOperation(path.Post)
		fixOperation(path.Put)
		fixOperation(path.Patch)
		fixOperation(path.Delete)
	}
}

func PrintSchema(s *fuego.Server) {
	FixSpec(s)
	spec := s.OutputOpenAPISpec()

	// Create ordered map for top-level keys
	orderedSpec := yaml.MapSlice{
		{Key: "info", Value: spec.Info},
		{Key: "openapi", Value: spec.OpenAPI},
		{Key: "components", Value: spec.Components},
		{Key: "paths", Value: spec.Paths},
		// {Key: "security", Value: spec.Security},
		// {Key: "servers", Value: spec.Servers},
		{Key: "tags", Value: spec.Tags},
		// {Key: "externalDocs", Value: spec.ExternalDocs},
	}
	bytes, err := yaml.MarshalWithOptions(orderedSpec, yaml.OmitEmpty())
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
