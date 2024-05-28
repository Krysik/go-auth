package auth

import "gorm.io/gorm"

func GetRefreshToken(db *gorm.DB, refreshToken, accountId string) (*RefreshToken, error) {
	var rt = RefreshToken{
		Token:     refreshToken,
		AccountId: accountId,
	}
	result := db.First(&rt)

	if result.Error != nil {
		return nil, result.Error
	}
	return &rt, nil
}
