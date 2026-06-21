package requests

import (
	"github.com/The-Notorious-Marauders/sorting-hat/core/constants"
	"github.com/The-Notorious-Marauders/sorting-hat/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golibs-starter/golib/web/response"
)

type Request interface {
	Validate(*validator.Validate) error
}

func Serialize[T Request](ctx *gin.Context, validator *validator.Validate, request T) bool {
	if err := ctx.ShouldBindWith(request, binding.FormMultipart); err != nil {
		response.WriteError(
			ctx.Writer,
			*utils.MakeException(constants.ErrorCode_BadRequest_CouldNotBindMultipartForm, err),
		)
		ctx.Abort()
		return false
	}

	if err := request.Validate(validator); err != nil {
		response.WriteError(
			ctx.Writer,
			*utils.MakeException(constants.ErrorCode_BadRequest_InvalidRequest, err),
		)
		ctx.Abort()
		return false
	}

	return true
}
