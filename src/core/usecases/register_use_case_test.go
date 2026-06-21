package usecases

import (
	"context"
	"strings"
	"testing"

	"github.com/The-Notorious-Marauders/sorting-hat/adapter/models"
	"github.com/The-Notorious-Marauders/sorting-hat/adapter/repositories/mocks"
	"github.com/The-Notorious-Marauders/sorting-hat/core/constants"
	"github.com/The-Notorious-Marauders/sorting-hat/core/entities/requests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestRegister_Success(t *testing.T) {
	userRepository := mocks.NewMockUserRepository(t)
	userRepository.EXPECT().
		FindByUsername(mock.Anything, "alice").
		Return(nil, gorm.ErrRecordNotFound)
	userRepository.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(user *models.User) bool {
			if user.Username != "alice" {
				return false
			}
			return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(testPassword)) == nil
		})).
		Run(func(_ context.Context, user *models.User) {
			user.Model = gorm.Model{ID: 7}
		}).
		Return(nil)

	useCase := NewRegisterUseCase(userRepository)

	response, exception := useCase.Register(context.Background(), requests.RegisterRequest{
		Username: "alice",
		Password: testPassword,
	})

	require.Nil(t, exception)
	require.NotNil(t, response)
	assert.Equal(t, uint(7), response.ID)
	assert.Equal(t, "alice", response.Username)
}

func TestRegister_UsernameAlreadyExists(t *testing.T) {
	userRepository := mocks.NewMockUserRepository(t)
	userRepository.EXPECT().
		FindByUsername(mock.Anything, "alice").
		Return(&models.User{Model: gorm.Model{ID: 1}, Username: "alice"}, nil)

	useCase := NewRegisterUseCase(userRepository)

	response, exception := useCase.Register(context.Background(), requests.RegisterRequest{
		Username: "alice",
		Password: testPassword,
	})

	assert.Nil(t, response)
	require.NotNil(t, exception)
	assert.EqualValues(t, constants.ErrorCode_Conflict_UsernameAlreadyExists, (*exception).Code())
}

func TestRegister_FindByUsernameError(t *testing.T) {
	userRepository := mocks.NewMockUserRepository(t)
	userRepository.EXPECT().
		FindByUsername(mock.Anything, "alice").
		Return(nil, assert.AnError)

	useCase := NewRegisterUseCase(userRepository)

	response, exception := useCase.Register(context.Background(), requests.RegisterRequest{
		Username: "alice",
		Password: testPassword,
	})

	assert.Nil(t, response)
	require.NotNil(t, exception)
	assert.EqualValues(t, constants.ErrorCode_InternalError_CouldNotFindUser, (*exception).Code())
}

func TestRegister_PasswordTooLong(t *testing.T) {
	userRepository := mocks.NewMockUserRepository(t)
	userRepository.EXPECT().
		FindByUsername(mock.Anything, "alice").
		Return(nil, gorm.ErrRecordNotFound)

	useCase := NewRegisterUseCase(userRepository)

	response, exception := useCase.Register(context.Background(), requests.RegisterRequest{
		Username: "alice",
		Password: strings.Repeat("a", 73),
	})

	assert.Nil(t, response)
	require.NotNil(t, exception)
	assert.EqualValues(t, constants.ErrorCode_BadRequest_PasswordTooLong, (*exception).Code())
}

func TestRegister_CreateFails(t *testing.T) {
	userRepository := mocks.NewMockUserRepository(t)
	userRepository.EXPECT().
		FindByUsername(mock.Anything, "alice").
		Return(nil, gorm.ErrRecordNotFound)
	userRepository.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(assert.AnError)

	useCase := NewRegisterUseCase(userRepository)

	response, exception := useCase.Register(context.Background(), requests.RegisterRequest{
		Username: "alice",
		Password: testPassword,
	})

	assert.Nil(t, response)
	require.NotNil(t, exception)
	assert.EqualValues(t, constants.ErrorCode_InternalError_CouldNotCreateUser, (*exception).Code())
}
