package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/The-Notorious-Marauders/sorting-hat/core/constants"
	"github.com/The-Notorious-Marauders/sorting-hat/core/entities"
	"github.com/The-Notorious-Marauders/sorting-hat/core/entities/requests"
	"github.com/The-Notorious-Marauders/sorting-hat/core/properties"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRefreshTokenUseCase_InvalidPrivateKey(t *testing.T) {
	jwksProps := &properties.JwksProperties{PrivateKey: "not-a-pem-key", PublicKey: "irrelevant"}

	useCase, err := NewRefreshTokenUseCase(jwksProps)

	assert.Error(t, err)
	assert.Nil(t, useCase)
}

func TestNewRefreshTokenUseCase_InvalidPublicKey(t *testing.T) {
	jwksProps, _ := generateTestJwksProperties(t)
	jwksProps.PublicKey = "not-a-pem-key"

	useCase, err := NewRefreshTokenUseCase(jwksProps)

	assert.Error(t, err)
	assert.Nil(t, useCase)
}

func TestRefreshToken_Success(t *testing.T) {
	jwksProps, privateKey := generateTestJwksProperties(t)
	issuedAt := time.Now().Add(-time.Hour)
	refreshToken, err := signToken(privateKey, "42", 12345, entities.TokenTypeRefresh, issuedAt, refreshTokenTTL)
	require.NoError(t, err)

	useCase, err := NewRefreshTokenUseCase(jwksProps)
	require.NoError(t, err)

	response, exception := useCase.RefreshToken(context.Background(), requests.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})

	require.Nil(t, exception)
	require.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.Equal(t, uint64(accessTokenTTL.Seconds()), response.AccessTokenExpiredTimeInSeconds)
	assert.Equal(t, uint64(refreshTokenTTL.Seconds()), response.RefreshTokenExpiredTimeInSeconds)

	accessClaims := parseAndVerifyToken(t, response.AccessToken, &privateKey.PublicKey)
	assert.Equal(t, "42", accessClaims.UserID)
	assert.Equal(t, 12345, accessClaims.LastLoginAtInSeconds)
	assert.Equal(t, entities.TokenTypeAccess, accessClaims.TokenType)

	refreshClaims := parseAndVerifyToken(t, response.RefreshToken, &privateKey.PublicKey)
	assert.Equal(t, "42", refreshClaims.UserID)
	assert.Equal(t, 12345, refreshClaims.LastLoginAtInSeconds)
	assert.Equal(t, entities.TokenTypeRefresh, refreshClaims.TokenType)
}

func TestRefreshToken_RejectsAccessToken(t *testing.T) {
	jwksProps, privateKey := generateTestJwksProperties(t)
	accessToken, err := signToken(privateKey, "42", 12345, entities.TokenTypeAccess, time.Now(), accessTokenTTL)
	require.NoError(t, err)

	useCase, err := NewRefreshTokenUseCase(jwksProps)
	require.NoError(t, err)

	response, exception := useCase.RefreshToken(context.Background(), requests.RefreshTokenRequest{
		RefreshToken: accessToken,
	})

	assert.Nil(t, response)
	require.NotNil(t, exception)
	assert.EqualValues(t, constants.ErrorCode_Unauthorized_InvalidJWT, (*exception).Code())
}

func TestRefreshToken_MalformedToken(t *testing.T) {
	jwksProps, _ := generateTestJwksProperties(t)

	useCase, err := NewRefreshTokenUseCase(jwksProps)
	require.NoError(t, err)

	response, exception := useCase.RefreshToken(context.Background(), requests.RefreshTokenRequest{
		RefreshToken: "this-is-not-a-jwt",
	})

	assert.Nil(t, response)
	require.NotNil(t, exception)
	assert.EqualValues(t, constants.ErrorCode_Unauthorized_MalformedJWT, (*exception).Code())
}

func TestRefreshToken_ExpiredToken(t *testing.T) {
	jwksProps, privateKey := generateTestJwksProperties(t)
	expiredToken, err := signToken(privateKey, "42", 12345, entities.TokenTypeRefresh, time.Now().Add(-2*refreshTokenTTL), refreshTokenTTL)
	require.NoError(t, err)

	useCase, err := NewRefreshTokenUseCase(jwksProps)
	require.NoError(t, err)

	response, exception := useCase.RefreshToken(context.Background(), requests.RefreshTokenRequest{
		RefreshToken: expiredToken,
	})

	assert.Nil(t, response)
	require.NotNil(t, exception)
	assert.EqualValues(t, constants.ErrorCode_Unauthorized_InvalidJWT, (*exception).Code())
}

func TestRefreshToken_WrongSignature(t *testing.T) {
	jwksProps, _ := generateTestJwksProperties(t)
	_, otherPrivateKey := generateTestJwksProperties(t)
	refreshToken, err := signToken(otherPrivateKey, "42", 12345, entities.TokenTypeRefresh, time.Now(), refreshTokenTTL)
	require.NoError(t, err)

	useCase, err := NewRefreshTokenUseCase(jwksProps)
	require.NoError(t, err)

	response, exception := useCase.RefreshToken(context.Background(), requests.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})

	assert.Nil(t, response)
	require.NotNil(t, exception)
	assert.EqualValues(t, constants.ErrorCode_Unauthorized_InvalidJWT, (*exception).Code())
}
