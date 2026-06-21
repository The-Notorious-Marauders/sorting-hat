package routers

import (
	"time"

	"github.com/The-Notorious-Marauders/sorting-hat/rest/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib"
	golibgin "github.com/golibs-starter/golib-gin"
	"github.com/golibs-starter/golib/web/actuator"
	"go.uber.org/fx"
)

type RegisterRoutersIn struct {
	fx.In
	App      *golib.App
	Engine   *gin.Engine
	Actuator *actuator.Endpoint

	AuthController *controllers.AuthController
}

func RegisterHandlers(app *golib.App, engine *gin.Engine) {
	engine.Use(golibgin.InitContext())
	engine.Use(golibgin.WrapAll(app.Handlers())...)
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func RegisterGinRouters(p RegisterRoutersIn) {
	base := p.Engine.Group(p.App.Path())
	base.GET("/actuator/health", gin.WrapF(p.Actuator.Health))
	base.GET("/actuator/info", gin.WrapF(p.Actuator.Info))

	api := base.Group("/api")
	api.POST("/sign-in", p.AuthController.SignIn)
	api.POST("/register", p.AuthController.Register)
	api.POST("/refresh-token", p.AuthController.RefreshToken)
}
