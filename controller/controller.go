package controller

import "github.com/go-fuego/fuego"

func NoTrailer[T, B any](controller func(fuego.ContextWithBody[B]) (T, error)) func(fuego.ContextWithBody[B]) (T, error) {
	return func(ctx fuego.ContextWithBody[B]) (T, error) {
		ctx.Response().Header().Del("Trailer")
		return controller(ctx)
	}
}
