package usecases

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"strconv"
	"time"

	"github.com/The-Notorious-Marauders/sorting-hat/adapter/repositories"
	"github.com/The-Notorious-Marauders/sorting-hat/core/constants"
	"github.com/The-Notorious-Marauders/sorting-hat/core/entities"
	"github.com/The-Notorious-Marauders/sorting-hat/core/entities/requests"
	"github.com/The-Notorious-Marauders/sorting-hat/core/entities/responses"
	"github.com/The-Notorious-Marauders/sorting-hat/core/properties"
	"github.com/The-Notorious-Marauders/sorting-hat/core/utils"
	"github.com/golibs-starter/golib/exception"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SignInUseCase interface {
	SignIn(context.Context, requests.SignInRequest) (*responses.SignInResponse, *exception.Exception)
}

type signInUseCaseImpl struct {
	jwksProps      *properties.JwksProperties
	userRepository repositories.UserRepository
	privateKey     *ecdsa.PrivateKey
}

func (s *signInUseCaseImpl) SignIn(
	ctx context.Context,
	request requests.SignInRequest,
) (
	*responses.SignInResponse,
	*exception.Exception,
) {
	user, err := s.userRepository.FindByUsername(ctx, request.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.MakeException(constants.ErrorCode_Unauthorized_InvalidCredentials, "invalid username or password")
		}
		return nil, utils.MakeException(constants.ErrorCode_InternalError_CouldNotFindUser, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		return nil, utils.MakeException(constants.ErrorCode_Unauthorized_InvalidCredentials, "invalid username or password")
	}

	userID := strconv.FormatUint(uint64(user.ID), 10)
	issuedAt := time.Now()
	lastLoginAtInSeconds := int(issuedAt.Unix())

	accessToken, err := signToken(s.privateKey, userID, lastLoginAtInSeconds, entities.TokenTypeAccess, issuedAt, accessTokenTTL)
	if err != nil {
		return nil, utils.MakeException(constants.ErrorCode_InternalError_CouldNotMakeJWTSignedString, err)
	}

	refreshToken, err := signToken(s.privateKey, userID, lastLoginAtInSeconds, entities.TokenTypeRefresh, issuedAt, refreshTokenTTL)
	if err != nil {
		return nil, utils.MakeException(constants.ErrorCode_InternalError_CouldNotMakeJWTSignedString, err)
	}

	if err := s.userRepository.UpdateLastLoginAt(ctx, user.ID, issuedAt); err != nil {
		return nil, utils.MakeException(constants.ErrorCode_InternalError_CouldNotUpdateLastLogin, err)
	}

	return &responses.SignInResponse{
		AccessToken:                      accessToken,
		AccessTokenExpiredTimeInSeconds:  uint64(accessTokenTTL.Seconds()),
		RefreshToken:                     refreshToken,
		RefreshTokenExpiredTimeInSeconds: uint64(refreshTokenTTL.Seconds()),
	}, nil
}

func NewSignInUseCase(
	jwksProps *properties.JwksProperties,
	userRepository repositories.UserRepository,
) (SignInUseCase, error) {
	privateKey, err := parseECDSAPrivateKey(jwksProps.PrivateKey)
	if err != nil {
		return nil, err
	}

	return &signInUseCaseImpl{
		jwksProps:      jwksProps,
		userRepository: userRepository,
		privateKey:     privateKey,
	}, nil
}
