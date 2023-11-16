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
	CreatedAt time.Time `gorm:"autoCreateTime:milli"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

type NewAccount struct {
	FullName string
	Email string
	Password string
}

func CreateAccount(db *gorm.DB, account NewAccount) (*Account, error) {
	acc := Account{
		FullName: account.FullName,
		Email: account.Email,
		// TODO: hash the password
		Password: account.Password,
	}

	result := db.Create(&acc)

	if result.Error != nil {
		return nil, result.Error
	}

	return &acc, nil
}


