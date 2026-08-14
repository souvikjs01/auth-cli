package service

import (
	"errors"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/souvikjs01/auth-cli/internals/models"
)

var (
	ErrInvalidTOTP       = errors.New("invalid authentication code")
	ErrMFAAlreadyEnabled = errors.New("MFA is already enabled")
	ErrMFANotEnabled     = errors.New("MFA is not enabled")
)

type TOTPService struct {
	issuer string
}

func NewTOTPService(issuer string) *TOTPService {
	return &TOTPService{
		issuer: issuer,
	}
}

func (s *TOTPService) GenerateKey(
	username string,
) (*otp.Key, error) {
	return totp.Generate(
		totp.GenerateOpts{
			Issuer:      s.issuer,
			AccountName: username,
		},
	)
}

func (s *TOTPService) Verify(
	code string,
	secret string,
) bool {
	return totp.Validate(code, secret)
}

func (s *AuthService) EnableMFA(
	user *models.User,
	secret string,
) error {
	if user.MFAEnabled {
		return ErrMFAAlreadyEnabled
	}

	user.MFASecret = &secret
	user.MFAEnabled = true

	return s.userRepo.Update(user)
}

func (s *AuthService) DisableMFA(
	user *models.User,
) error {
	if !user.MFAEnabled {
		return ErrMFANotEnabled
	}

	user.MFAEnabled = false
	user.MFASecret = nil

	return s.userRepo.Update(user)
}
