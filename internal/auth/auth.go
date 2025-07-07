package auth

import (
	"auth_service/internal/storage"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	store  *storage.Store
	logger *logrus.Logger
}

func NewService(store *storage.Store, logger *logrus.Logger) *AuthService {
	return &AuthService{
		store:  store,
		logger: logger,
	}
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func (a *AuthService) CreateTokenPair(ctx context.Context, userID uuid.UUID, userAgent string, IPAddress string) (*TokenPair, error) {
	accessToken, err := a.createAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.createRefreshToken()
	if err != nil {
		return nil, err
	}

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	session := storage.RefreshSession{
		UserID:     userID,
		TokenHash:  string(hashedRefreshToken),
		UserAgent:  userAgent,
		IPAddress:  IPAddress,
		ExpiryDate: time.Now().Add(24 * time.Hour),
		IsRevoked:  false,
	}

	if err := a.store.CreateRefreshSession(ctx, &session); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *AuthService) ParseAccessToken(accessToken string) (uuid.UUID, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, &jwt.MapClaims{})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("ошибка приведения к MapClaims")
	}

	sub, ok := claims["sub"]
	if !ok {
		return uuid.Nil, errors.New("subject в токене не найден")
	}

	id, err := uuid.Parse(sub.(string))
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
