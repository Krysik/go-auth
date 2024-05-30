package server

import (
	"errors"
	"fmt"
	"os"
)

type ENV struct {
	JWT_SECRET string
}

func NewEnv() (*ENV, error) {
	jwtSecret, err := getEnv("JWT_SECRET")

	if err != nil {
		return nil, err
	}

	return &ENV{
		JWT_SECRET: jwtSecret,
	}, nil
}

func getEnv(key string) (string, error) {
	env, exists := os.LookupEnv(key)

	if !exists || env == "" {
		return "", errors.New(fmt.Sprintf("variable \"%v\" is missing", key))
	}

	return env, nil
}
