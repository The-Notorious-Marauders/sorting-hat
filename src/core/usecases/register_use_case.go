package usecases

import (
	"context"
	"errors"

	"github.com/The-Notorious-Marauders/sorting-hat/adapter/models"
	"github.com/The-Notorious-Marauders/sorting-hat/adapter/repositories"
	"github.com/The-Notorious-Marauders/sorting-hat/core/constants"
	"github.com/The-Notorious-Marauders/sorting-hat/core/entities/requests"
	"github.com/The-Notorious-Marauders/sorting-hat/core/entities/responses"
	"github.com/The-Notorious-Marauders/sorting-hat/core/utils"
	"github.com/golibs-starter/golib/exception"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterUseCase interface {
	Register(context.Context, requests.RegisterRequest) (*responses.RegisterResponse, *exception.Exception)
}

type registerUseCaseImpl struct {
	userRepository repositories.UserRepository
}

func (r *registerUseCaseImpl) Register(
	ctx context.Context,
	request requests.RegisterRequest,
) (
	*responses.RegisterResponse,
	*exception.Exception,
) {
	_, err := r.userRepository.FindByUsername(ctx, request.Username)
	if err == nil {
		return nil, utils.MakeException(constants.ErrorCode_Conflict_UsernameAlreadyExists, "username already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.MakeException(constants.ErrorCode_InternalError_CouldNotFindUser, err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return nil, utils.MakeException(constants.ErrorCode_BadRequest_PasswordTooLong, err)
		}
		return nil, utils.MakeException(constants.ErrorCode_InternalError_CouldNotHashPassword, err)
	}

	user := &models.User{
		Username:     request.Username,
		PasswordHash: string(passwordHash),
	}

	if err := r.userRepository.Create(ctx, user); err != nil {
		return nil, utils.MakeException(constants.ErrorCode_InternalError_CouldNotCreateUser, err)
	}

	return &responses.RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

func NewRegisterUseCase(userRepository repositories.UserRepository) RegisterUseCase {
	return &registerUseCaseImpl{userRepository: userRepository}
}
