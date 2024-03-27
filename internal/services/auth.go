package auth

import (
	"auth/internal/models"
	"context"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log      *slog.Logger
	tokenTTL time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

// New returns Auth service
func New(log *slog.Logger, storage Storage) *Auth {
	return &Auth{}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	panic("emplement me")
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("registering new user")
	// Generate hash from service
	pHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.String(err.Error()))
		return 0, err
	}

	id, err := 
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	panic("implement me")
}
