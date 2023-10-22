package gin

import (
	"fmt"
	"github.com/delfimarime/fx/config"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

const (
	ApplicationJSONType        = "application/json"
	ApplicationProblemType     = "application/problem"
	ApplicationProblemJSONType = ApplicationProblemType + "+json"
)
const module = "gin"

func New(c config.Config) fx.Option {
	if c.Server.Type != "" && c.Server.Type != "gin" {
		return fx.Module(module)
	}
	return fx.Module(module, fx.Provide(
		func() *gin.Engine {
			mode := c.Server.Mode
			if mode == "" {
				mode = gin.ReleaseMode
			}
			gin.SetMode(mode)
			e := gin.Default()
			e.Use(gin.Logger())
			e.Use(gin.Recovery())
			e.Use(TryAndCatchHandler())
			if c.Server.Accept == nil {
				e.Use(AcceptMiddleware(ApplicationJSONType, ApplicationProblemType, ApplicationProblemJSONType))
			} else {
				e.Use(AcceptMiddleware(c.Server.Accept...))
			}
			if c.Server.ContentType == nil {
				e.Use(ContentTypeMiddleware(ApplicationJSONType))
			} else {
				e.Use(ContentTypeMiddleware(c.Server.ContentType...))
			}
			e.Use(ginzap.RecoveryWithZap(zap.L(), true))
			e.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
			return e
		}), fx.Invoke(func(e *gin.Engine, f ...gin.HandlerFunc) {
		for _, each := range f {
			e.Use(each)
		}
	}),
	)
}

func StartFxGin(c config.Config) fx.Option {
	return fx.Invoke(func(e *gin.Engine) error {
		if c.Server.Type != "" && c.Server.Type != "gin" {
			return nil
		}
		return e.Run(fmt.Sprintf(":%d", c.Server.Port))
	})
}
