package auth

import (
	"crypto/subtle"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        string `gorm:"primaryKey"`
	Token     string
	AccountId string
	CreatedAt time.Time
}

func ValidateCredentials(db *gorm.DB, email string, password string) (*Account, error) {
	var account Account

	queryResult := db.Where(&Account{Email: email}).First(&account)

	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	passwordMatch := comparePasswords(account.Password, account.Salt, password)

	if !passwordMatch {
		return nil, errors.New("invalid password")
	}

	return &account, nil
}

func comparePasswords(hashedPassword, salt, plainPassword string) bool {
	otherHash := hashPassword(plainPassword, salt)
	return subtle.ConstantTimeCompare([]byte(hashedPassword), otherHash) == 1
}

type AuthToken struct {
	AccessToken     string
	AccessTokenTtl  time.Time
	RefreshToken    string
	RefreshTokenTtl time.Time
}

type TokenClaims struct {
	jwt.RegisteredClaims
}

type TokenOpts struct {
	Issuer    string
	JwtSecret string
	Subject   string
}

func GenerateAuthTokens(opts TokenOpts) (*AuthToken, error) {
	secret := []byte(opts.JwtSecret)
	now := time.Now()

	accessTokenTtl := now.Add(time.Hour)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    opts.Issuer,
			ExpiresAt: jwt.NewNumericDate(accessTokenTtl),
			Subject:   opts.Subject,
		},
	})

	signedAccessJwt, accessTokenErr := accessToken.SignedString(secret)

	if accessTokenErr != nil {
		return nil, accessTokenErr
	}
	dayInHours := 24 * time.Hour

	refreshTokenTtl := now.Add(dayInHours * 30)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    opts.Issuer,
			ExpiresAt: jwt.NewNumericDate(refreshTokenTtl),
			Subject:   opts.Subject,
		},
	})
	signedRefreshJwt, refreshTokenErr := refreshToken.SignedString(secret)

	if refreshTokenErr != nil {
		return nil, refreshTokenErr
	}

	return &AuthToken{
		AccessToken:     signedAccessJwt,
		AccessTokenTtl:  accessTokenTtl,
		RefreshToken:    signedRefreshJwt,
		RefreshTokenTtl: refreshTokenTtl,
	}, nil
}

func SaveRefreshToken(db *gorm.DB, refreshToken, accountId string) error {
	rt := &RefreshToken{
		ID:        uuid.NewString(),
		Token:     refreshToken,
		AccountId: accountId,
	}

	createResult := db.Create(rt)

	if createResult.Error != nil {
		return createResult.Error
	}
	return nil
}
