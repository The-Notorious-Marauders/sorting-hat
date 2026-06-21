package usecases

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/The-Notorious-Marauders/sorting-hat/adapter/models"
	"github.com/The-Notorious-Marauders/sorting-hat/adapter/repositories/mocks"
	"github.com/The-Notorious-Marauders/sorting-hat/core/constants"
	"github.com/The-Notorious-Marauders/sorting-hat/core/entities"
	"github.com/The-Notorious-Marauders/sorting-hat/core/entities/requests"
	"github.com/The-Notorious-Marauders/sorting-hat/core/properties"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const testPassword = "correct-password"

func generateTestJwksProperties(t *testing.T) (*properties.JwksProperties, *ecdsa.PrivateKey) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	require.NoError(t, err)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privateKeyBytes})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	require.NoError(t, err)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes})

	return &properties.JwksProperties{
		PrivateKey: string(privateKeyPEM),
		PublicKey:  string(publicKeyPEM),
	}, privateKey
}

func hashPassword(t *testing.T, password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	return string(hash)
}

func TestNewSignInUseCase_InvalidPrivateKey(t *testing.T) {
	jwksProps := &properties.JwksProperties{PrivateKey: "not-a-pem-key", PublicKey: "irrelevant"}

	useCase, err := NewSignInUseCase(jwksProps, mocks.NewMockUserRepository(t))

	assert.Error(t, err)
	assert.Nil(t, useCase)
}

func TestSignIn_Success(t *testing.T) {
	jwksProps, privateKey := generateTestJwksProperties(t)
	lastLoginAt := time.Now().Add(-24 * time.Hour).Truncate(time.Second)
	user := &models.User{
		Model:        gorm.Model{ID: 42},
		Username:     "alice",
		PasswordHash: hashPassword(t, testPassword),
		LastLoginAt:  &lastLoginAt,
	}

	userRepository := mocks.NewMockUserRepository(t)
	userRepository.EXPECT().
		FindByUsername(mock.Anything, "alice").
		Return(user, nil)
	userRepository.EXPECT().
		UpdateLastLoginAt(mock.Anything, uint(42), mock.AnythingOfType("time.Time")).
		Return(nil)

	useCase, err := NewSignInUseCase(jwksProps, userRepository)
	require.NoError(t, err)

	beforeSignIn := time.Now()
	response, exception := useCase.SignIn(context.Background(), requests.SignInRequest{
		Username: "alice",
		Password: testPassword,
	})
	afterSignIn := time.Now()

	require.Nil(t, exception)
	require.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.Equal(t, uint64(accessTokenTTL.Seconds()), response.AccessTokenExpiredTimeInSeconds)
	assert.Equal(t, uint64(refreshTokenTTL.Seconds()), response.RefreshTokenExpiredTimeInSeconds)

	claims := parseAndVerifyToken(t, response.AccessToken, &privateKey.PublicKey)
	assert.Equal(t, "42", claims.UserID)
	assert.GreaterOrEqual(t, claims.LastLoginAtInSeconds, int(beforeSignIn.Unix()))
	assert.LessOrEqual(t, claims.LastLoginAtInSeconds, int(afterSignIn.Unix()))
}

func TestSignIn_UpdateLastLoginAtFails(t *testing.T) {
	jwksProps, _ := generateTestJwksProperties(t)
	user := &models.User{
		Model:        gorm.Model{ID: 42},
		Username:     "alice",
		PasswordHash: hashPassword(t, testPassword),
	}

	userRepository := mocks.NewMockUserRepository(t)
	userRepository.EXPECT().
		FindByUsername(mock.Anything, "alice").
		Return(user, nil)
	userRepository.EXPECT().
		UpdateLastLoginAt(mock.Anything, uint(42), mock.AnythingOfType("time.Time")).
		Return(assert.AnError)

	useCase, err := NewSignInUseCase(jwksProps, userRepository)
	require.NoError(t, err)

	response, exception := useCase.SignIn(context.Background(), requests.SignInRequest{
		Username: "alice",
		Password: testPassword,
	})

	assert.Nil(t, response)
	require.NotNil(t, exception)
	assert.EqualValues(t, constants.ErrorCode_InternalError_CouldNotUpdateLastLogin, (*exception).Code())
}

func TestSignIn_UserNotFound(t *testing.T) {
	jwksProps, _ := generateTestJwksProperties(t)

	userRepository := mocks.NewMockUserRepository(t)
	userRepository.EXPECT().
		FindByUsername(mock.Anything, "ghost").
		Return(nil, gorm.ErrRecordNotFound)

	useCase, err := NewSignInUseCase(jwksProps, userRepository)
	require.NoError(t, err)

	response, exception := useCase.SignIn(context.Background(), requests.SignInRequest{
		Username: "ghost",
		Password: testPassword,
	})

	assert.Nil(t, response)
	require.NotNil(t, exception)
	assert.EqualValues(t, constants.ErrorCode_Unauthorized_InvalidCredentials, (*exception).Code())
}

func TestSignIn_WrongPassword(t *testing.T) {
	jwksProps, _ := generateTestJwksProperties(t)
	user := &models.User{
		Model:        gorm.Model{ID: 1},
		Username:     "alice",
		PasswordHash: hashPassword(t, testPassword),
	}

	userRepository := mocks.NewMockUserRepository(t)
	userRepository.EXPECT().
		FindByUsername(mock.Anything, "alice").
		Return(user, nil)

	useCase, err := NewSignInUseCase(jwksProps, userRepository)
	require.NoError(t, err)

	response, exception := useCase.SignIn(context.Background(), requests.SignInRequest{
		Username: "alice",
		Password: "wrong-password",
	})

	assert.Nil(t, response)
	require.NotNil(t, exception)
	assert.EqualValues(t, constants.ErrorCode_Unauthorized_InvalidCredentials, (*exception).Code())
}

func TestSignIn_RepositoryError(t *testing.T) {
	jwksProps, _ := generateTestJwksProperties(t)

	userRepository := mocks.NewMockUserRepository(t)
	userRepository.EXPECT().
		FindByUsername(mock.Anything, "alice").
		Return(nil, assert.AnError)

	useCase, err := NewSignInUseCase(jwksProps, userRepository)
	require.NoError(t, err)

	response, exception := useCase.SignIn(context.Background(), requests.SignInRequest{
		Username: "alice",
		Password: testPassword,
	})

	assert.Nil(t, response)
	require.NotNil(t, exception)
	assert.EqualValues(t, constants.ErrorCode_InternalError_CouldNotFindUser, (*exception).Code())
}

func parseAndVerifyToken(t *testing.T, tokenString string, publicKey *ecdsa.PublicKey) *entities.SortingHatClaims {
	claims := &entities.SortingHatClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	require.NoError(t, err)
	return claims
}
