package auth

import "gorm.io/gorm"

func GetAccountById(db *gorm.DB, accountId string) (*Account, error) {
	var account = Account{
		ID: accountId,
	}

	result := db.First(&account)

	if result.Error != nil {
		return nil, result.Error
	}

	return &account, nil
}
