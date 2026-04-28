package session

import (
	"context"

	"github.com/alexedwards/scs/v2"
)

type Session struct {
	SessionManager *scs.SessionManager
	ctx            context.Context
}

func (s *Session) Exists(key string) bool {
	return s.SessionManager.Exists(s.ctx, key)
}
func (s *Session) Get(key string) interface{} {
	return s.SessionManager.Get(s.ctx, key)
}
func (s *Session) Pop(key string) interface{} {
	return s.SessionManager.Pop(s.ctx, key)
}
func (s *Session) Put(key string, val interface{}) {
	s.SessionManager.Put(s.ctx, key, val)
}
func (s *Session) RenewToken() error {
	return s.SessionManager.RenewToken(s.ctx)
}
func (s *Session) Destroy() error {
	return s.SessionManager.Destroy(s.ctx)
}

func GetSession(ctx context.Context) *Session {
	return &Session{
		SessionManager: GetSessionManager(),
		ctx:            ctx,
	}
}
