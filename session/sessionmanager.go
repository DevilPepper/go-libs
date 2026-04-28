package session

import (
	"net/http"
	"sync"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"

	"github.com/DevilPepper/go-libs/environment"
	"github.com/DevilPepper/go-libs/redis"
)

var (
	sessionManager *scs.SessionManager
	once           sync.Once
)

func GetSessionManager() *scs.SessionManager {
	if sessionManager == nil {
		once.Do(func() {
			sessionManager = scs.New()
			sessionManager.Lifetime = 24 * time.Hour
			sessionManager.IdleTimeout = 5 * time.Minute
			sessionManager.Cookie.HttpOnly = true
			sessionManager.Cookie.Secure = !environment.IsDev()
			sessionManager.Cookie.SameSite = http.SameSiteStrictMode
			sessionManager.Store = redisstore.New(redis.GetRedisPool())
		})
	}
	return sessionManager
}
