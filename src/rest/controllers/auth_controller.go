package controllers

import (
	"github.com/The-Notorious-Marauders/sorting-hat/core/entities/requests"
	"github.com/The-Notorious-Marauders/sorting-hat/core/usecases"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golibs-starter/golib/web/response"
)

type AuthController struct {
	signInUseCase       usecases.SignInUseCase
	registerUseCase     usecases.RegisterUseCase
	refreshTokenUseCase usecases.RefreshTokenUseCase
	validator           *validator.Validate
}

func NewAuthController(
	signInUseCase usecases.SignInUseCase,
	registerUseCase usecases.RegisterUseCase,
	refreshTokenUseCase usecases.RefreshTokenUseCase,
	validator *validator.Validate,
) *AuthController {
	return &AuthController{
		signInUseCase:       signInUseCase,
		registerUseCase:     registerUseCase,
		refreshTokenUseCase: refreshTokenUseCase,
		validator:           validator,
	}
}

func (a *AuthController) SignIn(ctx *gin.Context) {
	request := &requests.SignInRequest{}
	if !requests.Serialize(ctx, a.validator, request) {
		return
	}

	result, exception := a.signInUseCase.SignIn(ctx, *request)
	if exception != nil {
		response.WriteError(ctx.Writer, *exception)
		return
	}

	response.Write(ctx.Writer, response.Ok(result))
}

func (a *AuthController) Register(ctx *gin.Context) {
	request := &requests.RegisterRequest{}
	if !requests.Serialize(ctx, a.validator, request) {
		return
	}

	result, exception := a.registerUseCase.Register(ctx, *request)
	if exception != nil {
		response.WriteError(ctx.Writer, *exception)
		return
	}

	response.Write(ctx.Writer, response.Created(result))
}

func (a *AuthController) RefreshToken(ctx *gin.Context) {
	request := &requests.RefreshTokenRequest{}
	if !requests.Serialize(ctx, a.validator, request) {
		return
	}

	result, exception := a.refreshTokenUseCase.RefreshToken(ctx, *request)
	if exception != nil {
		response.WriteError(ctx.Writer, *exception)
		return
	}

	response.Write(ctx.Writer, response.Ok(result))
}
