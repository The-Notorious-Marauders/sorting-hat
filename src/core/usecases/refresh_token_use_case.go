package usecases

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"time"

	"github.com/The-Notorious-Marauders/sorting-hat/core/constants"
	"github.com/The-Notorious-Marauders/sorting-hat/core/entities"
	"github.com/The-Notorious-Marauders/sorting-hat/core/entities/requests"
	"github.com/The-Notorious-Marauders/sorting-hat/core/entities/responses"
	"github.com/The-Notorious-Marauders/sorting-hat/core/properties"
	"github.com/The-Notorious-Marauders/sorting-hat/core/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golibs-starter/golib/exception"
)

type RefreshTokenUseCase interface {
	RefreshToken(context.Context, requests.RefreshTokenRequest) (*responses.SignInResponse, *exception.Exception)
}

type refreshTokenUseCaseImpl struct {
	jwksProps  *properties.JwksProperties
	publicKey  *ecdsa.PublicKey
	privateKey *ecdsa.PrivateKey
}

func (r *refreshTokenUseCaseImpl) RefreshToken(
	ctx context.Context,
	request requests.RefreshTokenRequest,
) (
	*responses.SignInResponse,
	*exception.Exception,
) {
	claims := &entities.SortingHatClaims{}
	_, err := jwt.ParseWithClaims(request.RefreshToken, claims, func(token *jwt.Token) (any, error) {
		return r.publicKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, utils.MakeException(constants.ErrorCode_Unauthorized_MalformedJWT, err)
		}
		return nil, utils.MakeException(constants.ErrorCode_Unauthorized_InvalidJWT, err)
	}

	if claims.TokenType != entities.TokenTypeRefresh {
		return nil, utils.MakeException(constants.ErrorCode_Unauthorized_InvalidJWT, "token is not a refresh token")
	}

	issuedAt := time.Now()

	accessToken, err := signToken(r.privateKey, claims.UserID, claims.LastLoginAtInSeconds, entities.TokenTypeAccess, issuedAt, accessTokenTTL)
	if err != nil {
		return nil, utils.MakeException(constants.ErrorCode_InternalError_CouldNotMakeJWTSignedString, err)
	}

	refreshToken, err := signToken(r.privateKey, claims.UserID, claims.LastLoginAtInSeconds, entities.TokenTypeRefresh, issuedAt, refreshTokenTTL)
	if err != nil {
		return nil, utils.MakeException(constants.ErrorCode_InternalError_CouldNotMakeJWTSignedString, err)
	}

	return &responses.SignInResponse{
		AccessToken:                      accessToken,
		AccessTokenExpiredTimeInSeconds:  uint64(accessTokenTTL.Seconds()),
		RefreshToken:                     refreshToken,
		RefreshTokenExpiredTimeInSeconds: uint64(refreshTokenTTL.Seconds()),
	}, nil
}

func NewRefreshTokenUseCase(jwksProps *properties.JwksProperties) (RefreshTokenUseCase, error) {
	privateKey, err := parseECDSAPrivateKey(jwksProps.PrivateKey)
	if err != nil {
		return nil, err
	}

	publicKey, err := parseECDSAPublicKey(jwksProps.PublicKey)
	if err != nil {
		return nil, err
	}

	return &refreshTokenUseCaseImpl{
		jwksProps:  jwksProps,
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}
