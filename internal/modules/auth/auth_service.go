package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tamboto2000/otaqku-tasks/internal/common"
	"github.com/tamboto2000/otaqku-tasks/internal/config"
	"github.com/tamboto2000/otaqku-tasks/internal/dto"
)

var (
	ErrInvalidCredentials = common.Error{
		Code:    common.ErrCodeUnauthorized,
		Message: "Invalid credentials",
	}
	ErrInvalidToken = common.Error{
		Code:    common.ErrCodeUnauthorized,
		Message: "Invalid token",
	}
)

const (
	jwtScopeAccess  = "access"
	jwtScopeRefresh = "refresh"
)

type jwtClaims struct {
	jwt.RegisteredClaims
	Scope string `json:"scope"`
}

type AuthService struct {
	jwtCfg  config.JWT
	accRepo AccountRepository
	logger  *slog.Logger
}

func NewAuthService(jwtCfg config.JWT, accRepo AccountRepository, logger *slog.Logger) AuthService {
	return AuthService{jwtCfg: jwtCfg, accRepo: accRepo, logger: logger}
}

func (svc AuthService) RegisterAccount(ctx context.Context, req dto.CreateAccountRequest) error {
	acc, err := NewAccount(req)
	if err != nil {
		return err
	}

	exists, err := svc.accRepo.IsExistsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if exists {
		return common.Error{Code: common.ErrCodeAlreadyExists, Message: "account with the same email already exists"}
	}

	err = svc.accRepo.Save(ctx, acc)
	return err
}

func (svc AuthService) Login(ctx context.Context, email, pwd string) (dto.TokenResponse, error) {
	acc, err := svc.accRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return dto.TokenResponse{}, ErrInvalidCredentials
		}

		return dto.TokenResponse{}, err
	}

	if err := acc.MatchPassword(pwd); err != nil {
		return dto.TokenResponse{}, ErrInvalidCredentials
	}

	return svc.buildTokenPair(time.Now(), acc.ID)
}

func (svc AuthService) buildJwt(reqTime time.Time, d time.Duration, userId int, scope string) (dto.Token, error) {
	jti, err := uuid.NewV7()
	if err != nil {
		svc.logger.Error(fmt.Sprintf("error on generating token JTI: %v", err))
		return dto.Token{}, err
	}

	expiresAt := reqTime.Add(d)
	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID: jti.String(),
			// TODO: Add Issuer
			Subject:   strconv.Itoa(userId),
			IssuedAt:  jwt.NewNumericDate(reqTime),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(reqTime),
			// TODO: Add Audience
		},
		Scope: scope,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(svc.jwtCfg.SigningKey.Decoded)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("error on signing token: %v", err))
		return dto.Token{}, err
	}

	return dto.Token{
		Token:     tokenStr,
		ExpiresAt: expiresAt,
	}, nil
}

func (svc AuthService) ExchangeRefreshToken(ctx context.Context, tokenStr string) (dto.TokenResponse, error) {
	userId, err := svc.validateToken(tokenStr, jwtScopeRefresh)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	return svc.buildTokenPair(time.Now(), userId)
}

func (svc AuthService) buildTokenPair(reqTime time.Time, userId int) (dto.TokenResponse, error) {
	accessToken, err := svc.buildJwt(reqTime, time.Duration(svc.jwtCfg.AccessTokenDuration), userId, jwtScopeAccess)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	refreshToken, err := svc.buildJwt(reqTime, time.Duration(svc.jwtCfg.RefreshTokenDuration), userId, jwtScopeRefresh)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (svc AuthService) validateToken(tokenStr string, scope string) (int, error) {
	var claims jwtClaims

	_, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (any, error) {
		return svc.jwtCfg.SigningKey.Decoded, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed), errors.Is(err, jwt.ErrTokenUnverifiable),
			errors.Is(err, jwt.ErrTokenSignatureInvalid), errors.Is(err, jwt.ErrTokenExpired),
			errors.Is(err, jwt.ErrTokenUsedBeforeIssued):
			return 0, ErrInvalidToken
		}
	}

	if claims.Scope != scope {
		return 0, ErrInvalidToken
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("error on parsing subject as user id: %v", err))
		return 0, fmt.Errorf("parsing subject as user id error: %v", err)
	}

	return userId, nil
}

func (svc AuthService) ValidateAccessToken(ctx context.Context, tokenStr string) (int, error) {
	return svc.validateToken(tokenStr, jwtScopeAccess)
}
