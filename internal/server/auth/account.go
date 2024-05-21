package auth

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID	string `gorm:"primaryKey"`
	FullName string
	Email string
	Password string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewAccount struct {
	FullName string
	Email string
	Password string
}

func CreateAccount(db *gorm.DB, account NewAccount) (*Account, error) {
	acc := Account{
		ID: "1234",
		FullName: account.FullName,
		Email: account.Email,
		Password: account.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createResult := db.Create(&acc)

	if createResult.Error != nil {
		return nil, createResult.Error
	}

	return &acc, nil
}

func ListAccounts(db *gorm.DB) ([]Account, error) {
	var accounts []Account
	result := db.Find(&accounts)

	if result.Error != nil {
		return nil, result.Error
	}

	return accounts, nil
}


