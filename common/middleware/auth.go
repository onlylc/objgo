package middleware

import (
	"objgo/common/middleware/handler"
	"objgo/team/core/sdk/config"
	jwt "objgo/team/core/sdk/pkg/jwtauth"
	"time"
)

// AuthInit jwt验证new
func AuthInit() (*jwt.GinJWTMiddleware, error) {
	timeout := time.Hour
	if config.ApplicationConfig.Mode == "dev" {
		timeout = time.Duration(86010) * time.Hour
	} else {
		if config.JwtConfig.Timeout != 0 {
			timeout = time.Duration(config.JwtConfig.Timeout) * time.Second
		}
	}
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "test zone",
		Key:             []byte(config.JwtConfig.Secret),
		Timeout:         timeout,
		MaxRefresh:      time.Hour,
		PayloadFunc:     handler.PayloadFunc,
		IdentityHandler: handler.IdentityHandler,
		Authenticator:   handler.Authenticator,
		Authorizator:    handler.Authorizator,
		Unauthorized:    handler.Unauthorized,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})
}
