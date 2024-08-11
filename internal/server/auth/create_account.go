package auth

import (
	"encoding/base64"
	"time"

	"crypto/rand"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

type Account struct {
	ID        string `gorm:"primaryKey"`
	FullName  string
	Email     string
	Password  string
	Salt      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewAccount struct {
	FullName string
	Email    string
	Password string
}

func CreateAccount(db *gorm.DB, account NewAccount) (*Account, error) {
	salt, err := generatePasswordSalt(saltLength)

	if err != nil {
		return nil, err
	}

	hashedPassword := hashPassword(account.Password, salt)
	acc := &Account{
		ID:        uuid.NewString(),
		FullName:  account.FullName,
		Email:     account.Email,
		Password:  string(hashedPassword),
		Salt:      salt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createResult := db.Create(acc)

	if createResult.Error != nil {
		return nil, createResult.Error
	}

	return acc, nil
}

func ListAccounts(db *gorm.DB) ([]Account, error) {
	var accounts []Account
	result := db.Find(&accounts)

	if result.Error != nil {
		return nil, result.Error
	}

	return accounts, nil
}

const saltLength = 16
const hashTime = 2
const hashMemory = 64 * 1024
const hashThreads = 2
const keyLength = 32

func hashPassword(plainPassword string, salt string) []byte {
	hashedPassword := argon2.IDKey(
		[]byte(plainPassword),
		[]byte(salt),
		hashTime,
		hashMemory,
		hashThreads,
		keyLength,
	)

	return hashedPassword
}

func generatePasswordSalt(length int) (string, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)

	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}
