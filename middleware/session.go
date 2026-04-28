package middleware

import "github.com/DevilPepper/go-libs/session"

func SessionMiddleware() Middleware {
	return session.GetSessionManager().LoadAndSave
}
