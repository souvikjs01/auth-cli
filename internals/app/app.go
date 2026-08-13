package app

import (
	"github.com/souvikjs01/auth-cli/internals/models"
	"github.com/souvikjs01/auth-cli/internals/service"
)

type App struct {
	AuthService    *service.AuthService
	SessionService *service.SessionService

	CurrentSessionID string
}

func NewApp(
	authService *service.AuthService,
	sessionService *service.SessionService,
) *App {
	return &App{
		AuthService:    authService,
		SessionService: sessionService,
	}
}

func (a *App) Login(
	email string,
	password string,
) (*models.User, *models.Session, error) {
	user, err := a.AuthService.Login(
		service.LoginInput{
			Email:    email,
			Password: password,
		},
	)

	if err != nil {
		return nil, nil, err
	}

	session, err := a.SessionService.CreateSession(user.ID)
	if err != nil {
		return nil, nil, err
	}

	a.CurrentSessionID = session.ID

	return user, session, nil
}

func (a *App) CurrentUser() (*models.User, error) {
	if a.CurrentSessionID == "" {
		return nil, service.ErrNotAuthenticated
	}

	session, err := a.SessionService.GetValidSession(
		a.CurrentSessionID,
	)

	if err != nil {
		a.CurrentSessionID = ""
		return nil, err
	}

	return &session.User, nil
}

func (a *App) Logout() error {
	if a.CurrentSessionID == "" {
		return service.ErrNotAuthenticated
	}

	err := a.SessionService.Logout(
		a.CurrentSessionID,
	)

	if err != nil {
		return err
	}

	a.CurrentSessionID = ""

	return nil
}
