package middleware

import (
	"context"

	"github.com/bitdance-panic/gobuy/app/services/gateway/conf"
	"github.com/cloudwego/hertz/pkg/app"
)

func WhiteListMiddleware() app.HandlerFunc {
	whiteList := conf.GetConf().Auth.WhiteList
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Request.URI().Path())
		for _, p := range whiteList {
			if path == p {
				c.Set("skip_auth", true)
				break
			}
		}
		c.Next(ctx)
	}
}
