package config

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Env struct {
	stageRaw  string `ignored:"true"`
	Stage     Stage  `ignored:"true"`
	FrontUrl  string `env:"FRONT_URL"`
	JWTSecret string `env:"JWT_SECRET"`

	Log     Log     `env:", prefix=LOG_"`
	Service Service `env:", prefix=SERVICE_"`
	Vapid   Vapid   `env:", prefix=VAPID_"`
	Github  Github  `env:", prefix=GITHUB_"`

	DB DB `env:", prefix=DB_"`
}

type DB struct {
	URL string `env:"URL"`
}
type Log struct {
	Level string `env:"LEVEL"`
}

type Service struct {
	Port int `env:"PORT,default=8282"`
}

type Vapid struct {
	PublicKey  string `env:"PUBLIC_KEY"`
	PrivateKey string `env:"PRIVATE_KEY"`
}
type Github struct {
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
}

func NewEnv() (Env, error) {
	var env Env

	envPath := os.Getenv("CONFPATH")
	if envPath != "" {
		_ = godotenv.Load(envPath)
	}

	if err := envconfig.Process(context.Background(), &env); err != nil {
		return env, err
	}
	env.stageRaw = os.Getenv("STAGE")
	env.Stage = parseStage(env.stageRaw)

	if err := validateEnv(&env); err != nil {
		return env, err
	}

	return env, nil
}

func validateEnv(e *Env) error {
	if e.Stage == StageUnknown {
		return errUnknownStage
	}

	return nil

}
