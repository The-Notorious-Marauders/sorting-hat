package bootstrap

import (
	"github.com/The-Notorious-Marauders/sorting-hat/adapter/repositories"
	"github.com/The-Notorious-Marauders/sorting-hat/core/properties"
	"github.com/The-Notorious-Marauders/sorting-hat/core/usecases"
	"github.com/The-Notorious-Marauders/sorting-hat/rest/controllers"
	"github.com/The-Notorious-Marauders/sorting-hat/rest/routers"
	"github.com/go-playground/validator/v10"
	"github.com/golibs-starter/golib"
	golibdata "github.com/golibs-starter/golib-data"
	golibgin "github.com/golibs-starter/golib-gin"
	golibsec "github.com/golibs-starter/golib-security"
	"go.uber.org/fx"
)

func All() fx.Option {
	return fx.Options(
		golib.AppOpt(),
		golib.PropertiesOpt(),
		golib.LoggingOpt(),
		golib.EventOpt(),
		golibgin.GinHttpServerOpt(),
		golib.ActuatorEndpointOpt(),
		golib.BuildInfoOpt(Version, CommitHash, BuildTime),
		golib.HttpClientOpt(),
		golibsec.SecuredHttpClientOpt(),
		golibdata.DatasourceOpt(),

		fx.Provide(validator.New),

		golib.ProvideProps(properties.NewJwksProperties),

		fx.Provide(repositories.NewUserRepository),

		fx.Provide(usecases.NewSignInUseCase),
		fx.Provide(usecases.NewRegisterUseCase),
		fx.Provide(usecases.NewRefreshTokenUseCase),

		fx.Provide(controllers.NewAuthController),

		fx.Invoke(routers.RegisterHandlers),
		fx.Invoke(routers.RegisterGinRouters),

		golib.OnStopEventOpt(),
		golibgin.OnStopHttpServerOpt(),
	)
}
