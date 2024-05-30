package server

import (
	"github.com/caarlos0/env/v11"
)

type ENV struct {
	JwtSecret   string `env:"JWT_SECRET"`
	TokenIssuer string `env:"TOKEN_ISSUER"`
	Port        int    `env:"PORT" envDefault:"8080"`
	DbFilePath  string `env:"DB_FILE_PATH"`
}

func NewEnv() (*ENV, error) {
	e := &ENV{}

	if err := env.Parse(e); err != nil {
		return nil, err
	}
	return e, nil
}
