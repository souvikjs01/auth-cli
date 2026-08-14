package app

import (
	"github.com/chzyer/readline"
	"github.com/souvikjs01/auth-cli/internals/models"
	"github.com/souvikjs01/auth-cli/internals/service"
)

type App struct {
	AuthService    *service.AuthService
	SessionService *service.SessionService
	TokenService   *service.TOTPService

	CurrentSessionID string
	Readline         *readline.Instance
}

func NewApp(
	authService *service.AuthService,
	sessionService *service.SessionService,
	totpService *service.TOTPService,
) *App {
	return &App{
		AuthService:    authService,
		SessionService: sessionService,
		TokenService:   totpService,
	}
}

func (a *App) Login(
	username string,
	password string,
) (*models.User, *models.Session, error) {
	user, err := a.AuthService.Login(
		service.LoginInput{
			Username: username,
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

func (a *App) CurrentSession() (*models.Session, error) {
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

	return session, nil
}

func (a *App) CurrentUser() (*models.User, error) {
	session, err := a.CurrentSession()
	if err != nil {
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
