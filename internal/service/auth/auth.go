package auth

import (
	"banner/internal/domain"
	"banner/internal/pkg/jwt"
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserRepo interface {
	SaveUser(ctx context.Context, login string, passHash []byte) (err error)
	User(ctx context.Context, login string) (*domain.User, error)
}

type Auth struct {
	l        *zap.Logger
	r        UserRepo
	tokenTTL time.Duration
	secret   string
}

func New(
	l *zap.Logger,
	r UserRepo,
	tokenTTL time.Duration,
	secret string,
) *Auth {
	return &Auth{
		l:        l,
		r:        r,
		tokenTTL: tokenTTL,
		secret:   secret,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string) (string, error) {

	user, err := a.r.User(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			a.l.Warn("user not found", zap.Error(err))
			return "", domain.ErrInvalidCredentials
		}

		a.l.Error("failed to get user", zap.Error(err))
		return "", domain.ErrInternalServerError
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.l.Warn("invalid credentials", zap.Error(err))
		return "", domain.ErrInvalidCredentials
	}

	token, err := jwt.NewToken(user, a.tokenTTL, a.secret)
	if err != nil {
		return "", domain.ErrInternalServerError
	}

	a.l.Info("user logged in successfully")
	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, login string, password string) error {

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = a.r.SaveUser(ctx, login, passHash)
	if err != nil {
		if errors.Is(err, domain.ErrUserExists) {
			a.l.Warn("user already exists", zap.String("login", login))
			return err
		}
		a.l.Error("failed to save user", zap.Error(err))
		return err
	}

	a.l.Info("user registered")
	return nil
}

func (a *Auth) IsAdmin(ctx context.Context, token string) (bool, error) {
	cl, err := jwt.ParseToken(token, a.secret)
	if err != nil {
		return false, err
	}

	exp, err := cl.GetExpirationTime()
	if err != nil {
		return false, err
	}

	if exp.Before(time.Now()) {
		return false, domain.ErrInvalidToken
	}

	if cl["is_admin"].(bool) {
		return true, nil
	}
	return false, nil
}
