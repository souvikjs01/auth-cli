package service

import "errors"

var (
	ErrNotAuthenticated = errors.New("not authenticated")
	ErrSessionExpired   = errors.New("session has expired")
)
