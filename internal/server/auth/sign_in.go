package auth

import (
	"crypto/subtle"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func ValidateCredentials(db *gorm.DB, email string, password string) (*Account, error) {
	var account Account

	queryResult := db.Where(&Account{Email: email}).First(&account)

	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	passwordMatch, err := comparePasswords(account.Password, account.Salt, password)

	if err != nil {
		return nil, err
	}

	if !passwordMatch {
		return nil, errors.New("invalid password")
	}

	return &account, nil
}

func comparePasswords(hashedPassword, salt, plainPassword string) (bool, error) {
	otherHash := hashPassword(plainPassword, salt)

	if subtle.ConstantTimeCompare([]byte(hashedPassword), otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

type AuthToken struct {
	AccessToken     string
	AccessTokenTtl  time.Time
	RefreshToken    string
	RefreshTokenTtl time.Time
}

func GenerateAuthTokens(issuer string) (*AuthToken, error) {
	// TODO: get secret from environment variable
	secret := []byte("top-secret")
	now := time.Now()

	accessTokenTtl := now.Add(time.Hour)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": issuer,
		"exp": accessTokenTtl.Unix(),
	})
	signedAccessJwt, accessTokenErr := accessToken.SignedString(secret)

	if accessTokenErr != nil {
		return nil, accessTokenErr
	}
	dayInHours := 24 * time.Hour

	refreshTokenTtl := now.Add(dayInHours * 30)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": issuer,
		"exp": refreshTokenTtl.Unix(),
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
