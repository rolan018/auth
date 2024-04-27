package auth

import (
	"auth/internal/models"
	"auth/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"auth/internal/lib/jwt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotAdminRights     = errors.New("user doesn't have admin rights")
)

type Auth struct {
	log         *slog.Logger
	tokenTTL    time.Duration
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
}

type Storage interface {
	UserProvider
	UserSaver
	AppProvider
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
	CreateApp(ctx context.Context, app_name string, app_secret string) (int64, error)
}

// New returns Auth service
func New(log *slog.Logger, tokenTTL time.Duration, storage Storage) *Auth {

	return &Auth{log: log,
		tokenTTL:    tokenTTL,
		usrSaver:    storage,
		usrProvider: storage,
		appProvider: storage}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("attempting to login user")

	usr, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found")
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// check password
	err = bcrypt.CompareHashAndPassword(usr.PassHash, []byte(password))
	if err != nil {
		a.log.Info("invalid credentials", err)
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user logged in successfully")

	// generate new token
	token, err := jwt.NewToken(usr, app, a.tokenTTL)
	if err != nil {
		a.log.Info("failed to generate token", err)
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	return token, nil
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
		log.Error("failed to generate password hash", slog.String("op", err.Error()))
		return 0, err
	}

	id, err := a.usrSaver.SaveUser(ctx, email, pHash)
	if err != nil {
		log.Error("failed to save user", slog.String("op", err.Error()))
		return 0, err
	}
	log.Info("user registered")
	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)
	log.Info("checking if user is admin")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))
	return isAdmin, nil
}

func (a *Auth) CreateApp(ctx context.Context, email string, password string, app_name string, app_secret string) (int64, error) {
	const op = "auth.CreateApp"

	log := a.log.With(
		slog.String("op", op),
		slog.String("app_name", app_name),
	)
	// check profile
	usr, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found")
			return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// check password
	err = bcrypt.CompareHashAndPassword(usr.PassHash, []byte(password))
	if err != nil {
		a.log.Info("invalid credentials", err)
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	// check admin rights
	isAdmin, err := a.IsAdmin(ctx, usr.ID)
	if err != nil {
		return 0, err
	}
	if !isAdmin {
		log.Warn("user not a admin")
		return 0, fmt.Errorf("%s: %w", op, ErrNotAdminRights)
	}
	// create new app
	log.Info("registering new app")
	appID, err := a.appProvider.CreateApp(ctx, app_name, app_secret)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("app created")
	return appID, nil
}
