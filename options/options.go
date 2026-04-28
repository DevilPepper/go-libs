package options

import (
	"strings"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func OperationID(f any) func(*fuego.BaseRoute) {
	funcName := fuego.FuncName(f)
	parts := strings.Split(funcName, "/")
	operationID := parts[len(parts)-1]
	return func(r *fuego.BaseRoute) {
		r.Operation.OperationID = operationID
		r.FullName = funcName
	}
}

func NoTags() func(*fuego.BaseRoute) {
	return func(r *fuego.BaseRoute) {
		r.Operation.Tags = []string{}
	}
}

func NoContent() func(*fuego.BaseRoute) {
	return option.DefaultStatusCode(204)
}
