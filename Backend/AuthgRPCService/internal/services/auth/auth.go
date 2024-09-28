package auth

import (
	"authService/internal/domain/models"
	customErrors "authService/internal/errors"
	jwtTokken "authService/internal/lib/jwt"
	"context"
	"fmt"
	log "github.com/go-ozzo/ozzo-log"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"time"
	"unicode"
)

type Auth struct {
	logger       *log.Logger
	tokenTTL     time.Duration
	userSaver    UserSaver
	userProvider UserProvider
}

type UserSaver interface {
	SaveUser(ctx context.Context, user models.User) (string, error)
}

type UserProvider interface {
	FindUser(ctx context.Context, email string) (models.User, error)
}

func New(logger *log.Logger, tokenTTL time.Duration, saver UserSaver, provider UserProvider) *Auth {
	return &Auth{
		logger:       logger,
		tokenTTL:     tokenTTL,
		userSaver:    saver,
		userProvider: provider,
	}
}
func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (string, error) {
	const op = "Auth.RegisterNewUser"

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.logger.Error(fmt.Sprintf("%s hashing userModule", op), err.Error())
		return "", customErrors.UserAlreadyExists
	}

	if !validateNoSpaces(password) {
		a.logger.Error(fmt.Sprintf("%s validate password", op), "password contains spaces")
		return "", customErrors.InvalidInputData
	}

	user := models.User{
		Email:    email,
		PassHash: string(passHash),
	}

	validate := validator.New()

	err = validate.Struct(user)
	if err != nil {
		a.logger.Error(fmt.Sprintf("%s validate userModule", op), err.Error())
		return "", customErrors.InvalidInputData
	}

	_, err = a.userSaver.SaveUser(ctx, user)
	if err != nil {
		a.logger.Error(fmt.Sprintf("%s create userModule", op), err.Error())
		return "", err
	}

	token, err := jwtTokken.NewToken(user, a.tokenTTL)
	if err != nil {
		a.logger.Error(fmt.Sprintf("%s create token", op), err.Error())
		return "", err
	}

	return token, nil
}

func (a *Auth) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "Auth.Login"

	user, err := a.userProvider.FindUser(ctx, email)
	if err != nil {
		a.logger.Error(fmt.Sprintf("%s find userModule", op), err.Error())
		return "", customErrors.UserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password)); err != nil {
		a.logger.Error(fmt.Sprintf("%s compare password", op), err.Error())
		return "", err
	}

	token, err := jwtTokken.NewToken(user, a.tokenTTL)
	if err != nil {
		a.logger.Error(fmt.Sprintf("%s create token", op), err.Error())
		return "", err
	}
	return token, nil
}

func validateNoSpaces(input string) bool {
	for _, char := range input {
		if unicode.IsSpace(char) {
			return false
		}
	}
	return true
}
